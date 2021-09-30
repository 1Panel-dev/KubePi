package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"regexp"
	noesctmpl "text/template"
	"time"

	"github.com/KubeOperator/webkubectl/gotty/cache/token"
	"github.com/KubeOperator/webkubectl/gotty/pkg/randomstring"
	"github.com/KubeOperator/webkubectl/gotty/webtty"
	"github.com/NYTimes/gziphandler"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/KubeOperator/webkubectl/gotty/pkg/homedir"
)

// Server provides a webtty HTTP endpoint.
type Server struct {
	factory Factory
	options *Options

	upgrader         *websocket.Upgrader
	terminalTemplate *template.Template
	titleTemplate    *noesctmpl.Template
	cache            token.Cache
}


// New creates a new instance of Server.
// Server will use the New() of the factory provided to handle each request.
func New(factory Factory, options *Options, redisOptions *RedisOptions) (*Server, error) {

	indexData, err := Asset("static/terminal.html")
	if err != nil {
		panic("terminal not found") // must be in bindata
	}
	if options.IndexFile != "" {
		path := homedir.Expand(options.IndexFile)
		indexData, err = os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read custom index file at `%s`", path)
		}
	}
	terminalTemplate, err := template.New("index").Parse(string(indexData))
	if err != nil {
		panic("index template parse failed") // must be valid
	}

	titleTemplate, err := noesctmpl.New("title").Parse(options.TitleFormat)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse window title format `%s`", options.TitleFormat)
	}

	var originChecker func(r *http.Request) bool
	if options.WSOrigin != "" {
		matcher, err := regexp.Compile(options.WSOrigin)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compile regular expression of Websocket Origin: %s", options.WSOrigin)
		}
		originChecker = func(r *http.Request) bool {
			return matcher.MatchString(r.Header.Get("Origin"))
		}
	} else {
		// Don't check origin if ws-origin is not set
		originChecker = func(r *http.Request) bool {
			return true
		}
	}
	var cache token.Cache

	if redisOptions.UseRedisTokenCache == "true" {
		log.Println("Use redis to store token: ", redisOptions.Addr)
		client := redis.NewClient(redisOptions.Convert())
		cache = token.NewRedisCache(client, "kubeoperator-webkubectl-")
	} else {
		cache = token.NewMemCache()
	}

	return &Server{
		factory: factory,
		options: options,

		upgrader: &websocket.Upgrader{
			ReadBufferSize:  webtty.MaxBufferSize,
			WriteBufferSize: webtty.MaxBufferSize,
			Subprotocols:    webtty.Protocols,
			CheckOrigin:     originChecker,
		},
		terminalTemplate: terminalTemplate,
		titleTemplate:    titleTemplate,
		cache:            cache,
	}, nil
}

// Run starts the main process of the Server.
// The cancelation of ctx will shutdown the server immediately with aborting
// existing connections. Use WithGracefulContext() to support graceful shutdown.
func (server *Server) Run(ctx context.Context, options ...RunOption) error {
	cctx, cancel := context.WithCancel(ctx)
	opts := &RunOptions{gracefulCtx: context.Background()}
	for _, opt := range options {
		opt(opts)
	}

	counter := newCounter(time.Duration(server.options.Timeout) * time.Second)
	path := "/webkubectl/"
	customPath := os.Getenv("TERMINAL_PATH")
	if len(customPath) > 0 {
		path = customPath
	}

	if server.options.EnableRandomUrl {
		path = "/" + randomstring.Generate(server.options.RandomUrlLength) + "/"
	}

	handlers := server.setupHandlers(cctx, cancel, path, counter)
	srv, err := server.setupHTTPServer(handlers)
	if err != nil {
		return errors.Wrapf(err, "failed to setup an HTTP server")
	}

	if server.options.PermitWrite {
		log.Printf("Permitting clients to write input to the PTY.")
	}
	if server.options.Once {
		log.Printf("Once option is provided, accepting only one client")
	}

	if server.options.Port == "0" {
		log.Printf("Port number configured to `0`, choosing a random port")
	}
	hostPort := net.JoinHostPort(server.options.Address, server.options.Port)
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		return errors.Wrapf(err, "failed to listen at `%s`", hostPort)
	}

	scheme := "http"
	if server.options.EnableTLS {
		scheme = "https"
	}
	host, port, _ := net.SplitHostPort(listener.Addr().String())
	log.Printf("HTTP server is listening at: %s", scheme+"://"+host+":"+port)
	//if server.options.Address == "0.0.0.0" {
	//	for _, address := range listAddresses() {
	//		log.Printf("Alternative URL: %s", scheme+"://"+address+":"+port+path)
	//	}
	//}

	srvErr := make(chan error, 1)
	go func() {
		if server.options.EnableTLS {
			crtFile := homedir.Expand(server.options.TLSCrtFile)
			keyFile := homedir.Expand(server.options.TLSKeyFile)
			log.Printf("TLS crt file: " + crtFile)
			log.Printf("TLS key file: " + keyFile)

			err = srv.ServeTLS(listener, crtFile, keyFile)
		} else {
			err = srv.Serve(listener)
		}
		if err != nil {
			srvErr <- err
		}
	}()

	go func() {
		select {
		case <-opts.gracefulCtx.Done():
			srv.Shutdown(context.Background())
		case <-cctx.Done():
		}
	}()

	select {
	case err = <-srvErr:
		if err == http.ErrServerClosed { // by graceful ctx
			err = nil
		} else {
			cancel()
		}
	case <-cctx.Done():
		srv.Close()
		err = cctx.Err()
	}

	conn := counter.count()
	if conn > 0 {
		log.Printf("Waiting for %d connections to be closed", conn)
	}
	counter.wait()

	return err
}

func (server *Server) setupHandlers(ctx context.Context, cancel context.CancelFunc, pathPrefix string, counter *counter) http.Handler {
	staticFileHandler := http.FileServer(
		&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "static"},
	)

	var siteMux = http.NewServeMux()

	siteMux.HandleFunc(pathPrefix, server.handleIndex)
	siteMux.Handle(pathPrefix+"js/", http.StripPrefix(pathPrefix, staticFileHandler))
	siteMux.Handle(pathPrefix+"favicon.png", http.StripPrefix(pathPrefix, staticFileHandler))
	siteMux.Handle(pathPrefix+"css/", http.StripPrefix(pathPrefix, staticFileHandler))

	siteMux.HandleFunc(pathPrefix+"auth_token.js", server.handleAuthToken)
	siteMux.HandleFunc(pathPrefix+"config.js", server.handleConfig)
	siteMux.HandleFunc("/api/kube-config", server.handleKubeConfigApi)
	siteMux.HandleFunc("/api/kube-token", server.handleKubeTokenApi)
	if len(os.Getenv("TERMINAL_PATH")) < 1 {
		siteMux.HandleFunc("/", server.handleMain)
	}

	if os.Getenv("PPROF_ENABLED") != "" {
		siteMux.HandleFunc("/debug/pprof/", pprof.Index)
		siteMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		siteMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		siteMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		siteMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}
	siteHandler := http.Handler(siteMux)

	if server.options.EnableBasicAuth {
		log.Printf("Using Basic Authentication")
		siteHandler = server.wrapBasicAuth(siteHandler, server.options.Credential)
	}

	withGz := gziphandler.GzipHandler(server.wrapHeaders(siteHandler))
	siteHandler = server.wrapLogger(withGz)

	wsMux := http.NewServeMux()
	wsMux.Handle("/", siteHandler)
	wsMux.HandleFunc(pathPrefix+"ws", server.generateHandleWS(ctx, cancel, counter))
	siteHandler = http.Handler(wsMux)

	return siteHandler
}

func (server *Server) setupHTTPServer(handler http.Handler) (*http.Server, error) {
	srv := &http.Server{
		Handler: handler,
	}

	if server.options.EnableTLSClientAuth {
		tlsConfig, err := server.tlsConfig()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to setup TLS configuration")
		}
		srv.TLSConfig = tlsConfig
	}

	return srv, nil
}

func (server *Server) tlsConfig() (*tls.Config, error) {
	caFile := homedir.Expand(server.options.TLSCACrtFile)
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, errors.New("could not open CA crt file " + caFile)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("could not parse CA crt file data in " + caFile)
	}
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	return tlsConfig, nil
}
