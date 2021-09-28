package webkubectl

import "sync"

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
	return t.data[key]
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
