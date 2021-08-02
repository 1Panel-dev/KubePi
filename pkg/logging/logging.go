package logging

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

func GenLoggingSessionId() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	id := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(id, bytes)
	return string(id), nil
}

type LogSession struct {
	Id            string
	Bound         chan error
	sockJSSession sockjs.Session
}

type SessionMap struct {
	Sessions map[string]LogSession
	Lock     sync.Mutex
}

func (sm *SessionMap) Get(sessionId string) LogSession {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	return sm.Sessions[sessionId]
}

func (sm *SessionMap) Set(sessionId string, session LogSession) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	sm.Sessions[sessionId] = session
}

func (sm *SessionMap) Close(sessionId, reason string, status uint32) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	err := sm.Sessions[sessionId].sockJSSession.Close(status, reason)
	if err != nil {
		log.Println(err)
	}
	delete(sm.Sessions, sessionId)
}

var LogSessions = SessionMap{Sessions: make(map[string]LogSession)}

type LogMessage struct {
	SessionID string
	Data      string
}

func CreateLoggingHandler(path string) http.Handler {
	return sockjs.NewHandler(path, sockjs.DefaultOptions, logHandler)
}

func logHandler(session sockjs.Session) {

	var (
		buf        string
		err        error
		msg        LogMessage
		logSession LogSession
	)
	if buf, err = session.Recv(); err != nil {
		log.Printf("handleLogSession: can't Recv: %v", err)
		return
	}
	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		log.Printf("handleLogSession: can't UnMarshal (%v): %s", err, buf)
		return
	}
	if logSession = LogSessions.Get(msg.SessionID); logSession.Id == "" {
		log.Printf("handleLogSession: can't find session '%s'", msg.SessionID)
		return
	}
	logSession.sockJSSession = session
	LogSessions.Set(msg.SessionID, logSession)
	logSession.Bound <- nil
}

func WaitForLoggingStream(k8sClient kubernetes.Interface, namespace string, pod string, container string, tailLines int, follow bool, sessionId string) {
	select {
	case <-LogSessions.Get(sessionId).Bound:
		close(LogSessions.Get(sessionId).Bound)
		err := startLogProcess(k8sClient, namespace, pod, container, int64(tailLines), follow, LogSessions.Get(sessionId))
		if err != nil {
			LogSessions.Close(sessionId, err.Error(), 2)
			return
		}
		LogSessions.Close(sessionId, "Process exited", 1)
	}
}

func startLogProcess(k8sClient kubernetes.Interface, namespace string, pod string, container string, tailLines int64, follow bool, session LogSession) error {
	reader, err := k8sClient.CoreV1().
		Pods(namespace).
		GetLogs(pod, &v1.PodLogOptions{
			Container: container,
			Follow:    follow,
			TailLines: &tailLines,
		}).Stream(context.TODO())
	if err != nil {
		return err
	}

	ss := session.sockJSSession
	for {
		buf := make([]byte, 2048)
		numBytes, err := reader.Read(buf)
		if numBytes > 0 {
			message := string(buf[:numBytes])
			for _, val := range strings.Split(message, "\n") {
				if strings.TrimSpace(val) == "" {
					continue
				}
				err = ss.Send(strings.TrimSpace(val) + "\r\n")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		if err != nil {
			return err
		}
	}
}
