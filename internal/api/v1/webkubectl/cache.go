package webkubectl

import (
	"sync"
	"time"
)

type TerminalSessions struct {
	data  map[string]*Session
	mutex sync.Mutex
}

func NewTerminalSessions() *TerminalSessions {
	return &TerminalSessions{
		data: make(map[string]*Session),
	}
}

func (t *TerminalSessions) Get(key string) *Session {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	sess := t.data[key]
	if sess == nil {
		return nil
	}
	if !sess.ExpiresAt.IsZero() && time.Now().After(sess.ExpiresAt) {
		delete(t.data, key)
		return nil
	}
	return sess
}

func (t *TerminalSessions) Put(key string, sess *Session) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data[key] = sess
}

func (t *TerminalSessions) Delete(key string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	delete(t.data, key)
}
