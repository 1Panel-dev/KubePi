package shell

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"sync"
)

// 传输结束信号
const EndOfTransmission = "\0004"

type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

type TerminalSession struct {
	id            string
	bound         chan error
	sockJSSession sockjs.Session
	sizeChan      chan remotecommand.TerminalSize
	doneChan      chan struct{}
}

type TerminalMessage struct {
	Op, Data, SessionID string
	Rows, Cols          uint16
}

func (t TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

func (t TerminalSession) Read(p []byte) (int, error) {
	m, err := t.sockJSSession.Recv()
	if err != nil {
		return copy(p, EndOfTransmission), err
	}

	var msg TerminalMessage
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		return copy(p, EndOfTransmission), err
	}

	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
	default:
		return copy(p, EndOfTransmission), fmt.Errorf("unkown message type '%s'", msg.Op)
	}
	return len(p), nil
}

func (t TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Op:   "stdout",
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}
	if err = t.sockJSSession.Send(string(msg)); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (t TerminalSession) Toast(p string) error {
	msg, err := json.Marshal(TerminalMessage{
		Op:   "toast",
		Data: p,
	})
	if err != nil {
		return err
	}

	if err = t.sockJSSession.Send(string(msg)); err != nil {
		return err
	}

	return nil
}

type SessionMap struct {
	Sessions map[string]TerminalSession
	Lock     sync.RWMutex
}

func (sm *SessionMap) Get(sessionId string) TerminalSession {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	return sm.Sessions[sessionId]
}

func (sm *SessionMap) Set(sessionId string, session TerminalSession) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	sm.Sessions[sessionId] = session
}

func (sm *SessionMap) Close(sessionId string, status uint32, reason string) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	if err := sm.Sessions[sessionId].sockJSSession.Close(status, reason); err != nil {
		logrus.Error(err)
	}
	delete(sm.Sessions, sessionId)
}

var terminalSessions = SessionMap{Sessions: make(map[string]TerminalSession)}

func handleTerminalSession(session sockjs.Session) {
	var (
		buf             string
		err             error
		msg             TerminalMessage
		terminalSession TerminalSession
	)

	if buf, err = session.Recv(); err != nil {
		logrus.Errorf("handleTerminalSession: can't Recv: %v", err)
		return
	}
	if err := json.Unmarshal([]byte(buf), &msg); err != nil {
		logrus.Errorf("handleTerminalSession: can't UnMarshal: %v", err)
		return
	}
	if msg.Op != "bind" {
		logrus.Errorf("handleTerminalSession: expected 'bind' message, got: %s", buf)
		return
	}

	if terminalSession = terminalSessions.Get(msg.SessionID); terminalSession.id == "" {
		logrus.Errorf("handleTerminalSession: can't find session '%s'", msg.SessionID)
		return
	}
	terminalSession.sockJSSession = session
	terminalSessions.Set(msg.SessionID, terminalSession)
	terminalSession.bound <- nil
}

func CreateAttachHandler(path string) http.Handler {
	return sockjs.NewHandler(path, sockjs.DefaultOptions, handleTerminalSession)
}

func startProcess(k8sClient kubernetes.Interface, cfg *rest.Config, cmd []string, namespace string, podName string, containerName string, ptyHandler PtyHandler) error {

	req := k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: ptyHandler,
		Tty:               true,
	})

	if err != nil {
		return err
	}
	return nil
}

func GenTerminalSessionId() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	id := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(id, bytes)
	return string(id), nil
}

func isValidShell(validShells []string, shell string) bool {
	for _, validShell := range validShells {
		if validShell == shell {
			return true
		}
	}
	return false
}

// Waits for the SockJS connection to be opened by the client the session to be bound in handleTerminalSession
func WaitForTerminal(k8sClient kubernetes.Interface, cfg *rest.Config, namespace string, podName string, containerName string, sessionId string) {
	shell := "sh"

	select {
	case <-terminalSessions.Get(sessionId).bound:
		close(terminalSessions.Get(sessionId).bound)

		var err error
		validShells := []string{"bash", "sh", "powershell", "cmd"}

		if isValidShell(validShells, shell) {
			cmd := []string{shell}
			err = startProcess(k8sClient, cfg, cmd, namespace, podName, containerName, terminalSessions.Get(sessionId))
		} else {
			// No shell given or it was not valid: try some shells until one succeeds or all fail
			// FIXME: if the first shell fails then the first keyboard event is lost
			for _, testShell := range validShells {
				cmd := []string{testShell}
				if err = startProcess(k8sClient, cfg, cmd, namespace, podName, containerName, terminalSessions.Get(sessionId)); err == nil {
					break
				}
			}
		}

		if err != nil {
			terminalSessions.Close(sessionId, 2, err.Error())
			return
		}

		terminalSessions.Close(sessionId, 1, "Process exited")
	}
}
