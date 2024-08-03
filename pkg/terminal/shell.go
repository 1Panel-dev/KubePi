package terminal

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/utils/ptr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const END_OF_TRANSMISSION = "\u0004"
const SessionTerminalStoreTime = 5 // session timeout (minute)

// PtyHandler is what remotecommand expects from a pty
type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

// TerminalSession implements PtyHandler (using a SockJS connection)
type TerminalSession struct {
	Id            string
	Bound         chan error
	sockJSSession sockjs.Session
	SizeChan      chan remotecommand.TerminalSize
	doneChan      chan struct{}
	TimeOut       time.Time
}

// TerminalMessage is the messaging protocol between ShellController and TerminalSession.
//
// OP      DIRECTION  FIELD(S) USED  DESCRIPTION
// ---------------------------------------------------------------------
// bind    fe->be     SessionID      Id sent back from TerminalResponse
// stdin   fe->be     Data           Keystrokes/paste buffer
// resize  fe->be     Rows, Cols     New terminal size
// stdout  be->fe     Data           Output from the process
// toast   be->fe     Data           OOB message to be shown to the user
type TerminalMessage struct {
	Op, Data, SessionID string
	Rows, Cols          uint16
}

// TerminalSize handles pty->process resize events
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.SizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

// Read handles pty->process messages (stdin, resize)
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Read(p []byte) (int, error) {
	session := TerminalSessions.Get(t.Id)
	if session.TimeOut.Before(time.Now()) {
		_ = TerminalSessions.Sessions[session.Id].sockJSSession.Close(2, "the connection has been disconnected. Please reconnect")
		return 0, errors.New("the connection has been disconnected. Please reconnect")
	}
	TerminalSessions.Set(session.Id, session)
	m, err := session.sockJSSession.Recv()
	if err != nil {
		// Send terminated signal to process to avoid resource leak
		return copy(p, END_OF_TRANSMISSION), err
	}

	var msg TerminalMessage
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		return copy(p, END_OF_TRANSMISSION), err
	}

	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		session.SizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	default:
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("unknown message type '%s'", msg.Op)
	}
}

// Write handles process->pty stdout
// Called from remotecommand whenever there is any output
func (t TerminalSession) Write(p []byte) (int, error) {
	session := TerminalSessions.Get(t.Id)
	if session.TimeOut.Before(time.Now()) {
		_ = TerminalSessions.Sessions[session.Id].sockJSSession.Close(2, "the connection has been disconnected. Please reconnect")
		return 0, errors.New("the connection has been disconnected. Please reconnect")
	}
	TerminalSessions.Set(session.Id, session)
	msg, err := json.Marshal(TerminalMessage{
		Op:   "stdout",
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}

	if err = session.sockJSSession.Send(string(msg)); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Toast can be used to send the user any OOB messages
// hterm puts these in the center of the terminal
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

// SessionMap stores a map of all TerminalSession objects and a lock to avoid concurrent conflict
type SessionMap struct {
	Sessions map[string]TerminalSession
	Lock     sync.RWMutex
}

// Get return a given terminalSession by sessionId
func (sm *SessionMap) Get(sessionId string) TerminalSession {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	return sm.Sessions[sessionId]
}

// Set store a TerminalSession to SessionMap
func (sm *SessionMap) Set(sessionId string, session TerminalSession) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	session.TimeOut = time.Now().Add(SessionTerminalStoreTime * time.Minute)
	sm.Sessions[sessionId] = session
}

// Close shuts down the SockJS connection and sends the status code and reason to the client
// Can happen if the process exits or if there is an error starting up the process
// For now the status code is unused and reason is shown to the user (unless "")
func (sm *SessionMap) Close(sessionId string, status uint32, reason string) {
	if _, ok := sm.Sessions[sessionId]; !ok {
		return
	}
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	err := sm.Sessions[sessionId].sockJSSession.Close(status, reason)
	if err != nil && status != 1 {
		log.Println(err)
	}

	delete(sm.Sessions, sessionId)
}

// Clean all session when system logout
func (sm *SessionMap) Clean() {
	for _, v := range sm.Sessions {
		v.sockJSSession.Close(2, "system is logout, please retry...")
	}
	sm.Sessions = make(map[string]TerminalSession)
}

var TerminalSessions = SessionMap{Sessions: make(map[string]TerminalSession)}

// handleTerminalSession is Called by net/http for any new /api/sockjs connections
func handleTerminalSession(session sockjs.Session) {
	var (
		buf             string
		err             error
		msg             TerminalMessage
		terminalSession TerminalSession
	)

	if buf, err = session.Recv(); err != nil {
		log.Printf("handleTerminalSession: can't Recv: %v", err)
		return
	}

	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		log.Printf("handleTerminalSession: can't UnMarshal (%v): %s", err, buf)
		return
	}

	if msg.Op != "bind" {
		log.Printf("handleTerminalSession: expected 'bind' message, got: %s", buf)
		return
	}

	if terminalSession = TerminalSessions.Get(msg.SessionID); terminalSession.Id == "" {
		log.Printf("handleTerminalSession: can't find session '%s'", msg.SessionID)
		return
	}

	terminalSession.sockJSSession = session
	TerminalSessions.Set(msg.SessionID, terminalSession)
	terminalSession.Bound <- nil
}

// CreateAttachHandler is called from main for /api/sockjs
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
	ctx :=context.Background()
	err = exec.StreamWithContext(ctx,remotecommand.StreamOptions{
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

// isValidShell checks if the shell is an allowed one
func isValidShell(validShells []string, shell string) bool {
	for _, validShell := range validShells {
		if validShell == shell {
			return true
		}
	}
	return false
}

// WaitForTerminal is called from apihandler.handleAttach as a goroutine
// Waits for the SockJS connection to be opened by the client the session to be Bound in handleTerminalSession
func WaitForTerminal(k8sClient kubernetes.Interface, cfg *rest.Config, namespace string, podName string, containerName string, sessionId string, shell string) {
	select {
	case <-TerminalSessions.Get(sessionId).Bound:
		close(TerminalSessions.Get(sessionId).Bound)

		var err error
		validShells := []string{shell}

		if isValidShell(validShells, shell) {
			cmd := []string{shell}
			err = startProcess(k8sClient, cfg, cmd, namespace, podName, containerName, TerminalSessions.Get(sessionId))
		} else {
			// No shell given or it was not valid: try some shells until one succeeds or all fail
			// FIXME: if the first shell fails then the first keyboard event is lost
			for _, testShell := range validShells {
				cmd := []string{testShell}
				if err = startProcess(k8sClient, cfg, cmd, namespace, podName, containerName, TerminalSessions.Get(sessionId)); err == nil {
					break
				}
			}
		}

		if err != nil {
			TerminalSessions.Close(sessionId, 2, err.Error())
			return
		}

		TerminalSessions.Close(sessionId, 1, "Process exited")
	}
}

//node shell
func WaitForNodeShellTerminal(k8sClient kubernetes.Interface, cfg *rest.Config,nodeName string, sessionId string) {
	select {
	case <-TerminalSessions.Get(sessionId).Bound:
		close(TerminalSessions.Get(sessionId).Bound)

		err := startNodeShellProcess(k8sClient, cfg, nodeName, TerminalSessions.Get(sessionId))
		

		if err != nil {
			TerminalSessions.Close(sessionId, 2, err.Error())
			return
		}

		TerminalSessions.Close(sessionId, 1, "Process exited")
	}
}
func startNodeShellProcess(k8sClient kubernetes.Interface, cfg *rest.Config, nodeName string, ptyHandler PtyHandler) error {

	/*创建特权容器*/
	containerName := "node-shell"
	podName :="node-"+nodeName+"-shell"
	namespace :="default"
	image := "alpine:latest"
	ctx :=context.Background()

	getpod ,err :=k8sClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
	    FieldSelector :"metadata.name="+podName+",metadata.namespace="+namespace,
	})
	if err != nil {
	    return err
	}
	//没有pod则创建
	if len(getpod.Items) == 0 {
	    container := &v1.Container{
			Name:            containerName,
			Image:           image,
			ImagePullPolicy: v1.PullIfNotPresent,
			Stdin:           true,
			StdinOnce:       true,
			TTY:             true,
			Command:         []string{"/bin/sh"},
			SecurityContext: &v1.SecurityContext{
				Privileged : ptr.To(true),
			},
			//挂载hostPath到容器中
			VolumeMounts: []v1.VolumeMount{
			    {
					Name:      "host-root",
					MountPath: "/host",
				},
			},
		}
	
	
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: podName,
				Namespace: namespace,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{*container},
				NodeName: nodeName,
				RestartPolicy: v1.RestartPolicyNever,
				HostNetwork: true,
				HostPID: true,
				Tolerations: []v1.Toleration{
				  
					{
							Key: "CriticalAddonsOnly",
							Operator: v1.TolerationOpExists,
					},
					{
							Key: "NoExecute",
							Operator: v1.TolerationOpExists,
					},
				},
				//挂载hostPath到容器中
				Volumes: []v1.Volume{
				    v1.Volume{
						Name: "host-root",
						VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: "/",
								},
						},
					},
				},
			},
		}
		_, err = k8sClient.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{
		});
		if err != nil {
			return err
		}
		//等待pod创建完成
		time.Sleep(5 * time.Second)
	}

	

	cmd :=[]string{"sh"}
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
	
	err = exec.StreamWithContext(ctx,remotecommand.StreamOptions{
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