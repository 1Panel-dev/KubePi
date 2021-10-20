package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"

	cache "github.com/KubeOperator/webkubectl/gotty/cache/token"
	"github.com/KubeOperator/webkubectl/gotty/pkg/randomstring"
	"github.com/KubeOperator/webkubectl/gotty/webtty"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (server *Server) generateHandleWS(ctx context.Context, cancel context.CancelFunc, counter *counter) http.HandlerFunc {
	once := new(int64)

	go func() {
		select {
		case <-counter.timer().C:
			cancel()
		case <-ctx.Done():
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		if server.options.Once {
			success := atomic.CompareAndSwapInt64(once, 0, 1)
			if !success {
				http.Error(w, "Server is shutting down", http.StatusServiceUnavailable)
				return
			}
		}

		num := counter.add(1)
		closeReason := "unknown reason"

		defer func() {
			num := counter.done()
			log.Printf(
				"Connection closed: %s, reason: %s, connections: %d/%d",
				r.RemoteAddr, closeReason, num, server.options.MaxConnection,
			)

			if server.options.Once {
				cancel()
			}
		}()

		if int64(server.options.MaxConnection) != 0 {
			if num > server.options.MaxConnection {
				closeReason = "exceeding max number of connections"
				return
			}
		}

		log.Printf("New client connected: %s, connections: %d/%d", r.RemoteAddr, num, server.options.MaxConnection)

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		server.upgrader.ReadBufferSize = webtty.MaxBufferSize
		server.upgrader.WriteBufferSize = webtty.MaxBufferSize
		server.upgrader.EnableCompression = true
		conn, err := server.upgrader.Upgrade(w, r, nil)
		if err != nil {
			closeReason = err.Error()
			return
		}
		defer conn.Close()
		conn.SetCompressionLevel(9)
		err = server.processWSConn(ctx, conn)

		switch err {
		case ctx.Err():
			closeReason = "cancelation"
		case webtty.ErrSlaveClosed:
			closeReason = server.factory.Name()
		case webtty.ErrMasterClosed:
			closeReason = "client close"
		case webtty.ErrConnectionLostPing:
			closeReason = webtty.ErrConnectionLostPing.Error()
		default:
			closeReason = fmt.Sprintf("an error: %s", err)
		}
	}
}

func (server *Server) processWSConn(ctx context.Context, conn *websocket.Conn) error {
	typ, initLine, err := conn.ReadMessage()
	if err != nil {
		return errors.Wrapf(err, "failed to authenticate websocket connection")
	}
	if typ != websocket.TextMessage {
		return errors.New("failed to authenticate websocket connection: invalid message type")
	}

	var init InitMessage
	err = json.Unmarshal(initLine, &init)
	if err != nil {
		return errors.Wrapf(err, "failed to authenticate websocket connection")
	}
	if init.AuthToken != server.options.Credential {
		return errors.New("failed to authenticate websocket connection")
	}

	queryPath := "?"
	if server.options.PermitArguments && init.Arguments != "" {
		queryPath = init.Arguments
	}

	query, err := url.Parse(queryPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse arguments")
	}
	windowTitle := "Webkubectl"
	params := query.Query()
	params.Del("arg")
	arg := params.Get("token")
	params.Add("arg", arg)
	log.Println("arg: " + arg)
	var slave Slave
	slave, err = server.factory.New(params)
	if err != nil {
		return errors.Wrapf(err, "failed to create backend")
	}
	defer slave.Close()

	titleVars := server.titleVariables(
		[]string{"server", "master", "slave"},
		map[string]map[string]interface{}{
			"server": server.options.TitleVariables,
			"master": {
				"remote_addr": conn.RemoteAddr(),
			},
			"slave": slave.WindowTitleVariables(),
		},
	)

	titleBuf := new(bytes.Buffer)

	err = server.titleTemplate.Execute(titleBuf, titleVars)
	if err != nil {
		return errors.Wrapf(err, "failed to fill window title template")
	}
	if len(windowTitle) > 0 {
		titleBuf.Reset()
		titleBuf.WriteString(windowTitle)
	}
	opts := []webtty.Option{
		webtty.WithWindowTitle(titleBuf.Bytes()),
	}
	if server.options.PermitWrite {
		opts = append(opts, webtty.WithPermitWrite())
	}
	if server.options.EnableReconnect {
		opts = append(opts, webtty.WithReconnect(server.options.ReconnectTime))
	}
	if server.options.Width > 0 {
		opts = append(opts, webtty.WithFixedColumns(server.options.Width))
	}
	if server.options.Height > 0 {
		opts = append(opts, webtty.WithFixedRows(server.options.Height))
	}
	if server.options.Preferences != nil {
		opts = append(opts, webtty.WithMasterPreferences(server.options.Preferences))
	}

	tty, err := webtty.New(&wsWrapper{conn}, slave, opts...)
	if err != nil {
		return errors.Wrapf(err, "failed to create webtty")
	}

	err = tty.Run(ctx)

	return err
}

func (server *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	titleVars := server.titleVariables(
		[]string{"server", "master"},
		map[string]map[string]interface{}{
			"server": server.options.TitleVariables,
			"master": {
				"remote_addr": r.RemoteAddr,
			},
		},
	)

	titleBuf := new(bytes.Buffer)
	err := server.titleTemplate.Execute(titleBuf, titleVars)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	indexVars := map[string]interface{}{
		"title": titleBuf.String(),
	}

	indexBuf := new(bytes.Buffer)
	err = server.terminalTemplate.Execute(indexBuf, indexVars)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Write(indexBuf.Bytes())
}

func (server *Server) handleMain(w http.ResponseWriter, r *http.Request) {
	indexData, err := Asset("static/index.html")
	if err != nil {
		panic("index not found") // must be in bindata
	}
	indexData = bytes.Replace(indexData, []byte("#app.Version"), []byte(os.Getenv("app.Version")), -1)
	if len(os.Getenv("HIDE_INDEX_FOOTER")) > 0 {
		w.Write(bytes.Replace(indexData, []byte("footer-show"), []byte(""), -1))
	} else {
		w.Write(indexData)
	}

}

func (server *Server) handleAuthToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	// @TODO hashing?
	w.Write([]byte("var gotty_auth_token = '" + server.options.Credential + "';"))
}

func (server *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte("var gotty_term = '" + server.options.Term + "';"))
}

// titleVariables merges maps in a specified order.
// varUnits are name-keyed maps, whose names will be iterated using order.
func (server *Server) titleVariables(order []string, varUnits map[string]map[string]interface{}) map[string]interface{} {
	titleVars := map[string]interface{}{}

	for _, name := range order {
		vars, ok := varUnits[name]
		if !ok {
			panic("title variable name error")
		}
		for key, val := range vars {
			titleVars[key] = val
		}
	}

	// safe net for conflicted keys
	for _, name := range order {
		titleVars[name] = varUnits[name]
	}

	return titleVars
}

func (server *Server) handleKubeConfigApi(w http.ResponseWriter, r *http.Request) {
	result := ApiResponse{
		Success: false,
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if !strings.EqualFold(r.Method, "POST") {
		result.Message = "Method Not Allowed"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(405)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		result.Message = "Invalid Request Body"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(406)
		return
	}
	var request KubeConfigRequest
	if err = json.Unmarshal(body, &request); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		result.Message = "Invalid Request Body"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(406)
		return
	}
	if len(request.KubeConfig) < 10 {
		result.Message = "Invalid Kube Config"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(406)
		return
	}
	if _, base64err := base64.StdEncoding.DecodeString(request.KubeConfig); base64err != nil {
		result.Message = "Invalid Base64 Content"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(406)
		return
	}

	//fmt.Printf("%+v", requst)
	token := randomstring.Generate(20)
	ttyParameter := cache.TtyParameter{
		Title: request.Name,
		Arg:   strings.Replace(request.KubeConfig, " ", "", -1),
	}
	if err := server.cache.Add(token, &ttyParameter, time.Duration(server.options.TokenExpiresDuration)*time.Second); err != nil {
		log.Printf("save token and ttyParam err:%s", err.Error())
		result.Success = false
		result.Message = err.Error()
		json.NewEncoder(w).Encode(result)
		return
	}
	result.Success = true
	result.Token = token
	json.NewEncoder(w).Encode(result)
}

func (server *Server) handleKubeTokenApi(w http.ResponseWriter, r *http.Request) {
	result := ApiResponse{
		Success: false,
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if !strings.EqualFold(r.Method, "POST") {
		result.Message = "Method Not Allowed"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(405)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		result.Message = "Invalid Request Body"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(406)
		return
	}
	var request KubeTokenRequest
	if err = json.Unmarshal(body, &request); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		result.Message = "Invalid Request Body"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(406)
		return
	}
	if !strings.HasPrefix(request.ApiServer, "http") {
		result.Message = "Invalid ApiServer"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(406)
		return
	}

	if len(request.Token) < 10 {
		result.Message = "Invalid Bearer Token"
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(406)
		return
	}
	//fmt.Printf("%+v", requst)
	token := randomstring.Generate(20)
	ttyParameter := cache.TtyParameter{
		Title: request.Name,
		Arg:   strings.Replace(request.ApiServer, " ", "", -1) + " " + strings.Replace(request.Token, " ", "", -1),
	}
	if err := server.cache.Add(token, &ttyParameter, time.Duration(server.options.TokenExpiresDuration)*time.Second); err != nil {
		log.Printf("save token and ttyParam err:%s", err.Error())
		result.Success = false
		result.Message = err.Error()
		json.NewEncoder(w).Encode(result)
		return
	}
	result.Success = true
	result.Token = token
	json.NewEncoder(w).Encode(result)
}
