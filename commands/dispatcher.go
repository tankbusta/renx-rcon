package commands

import (
	"container/list"
	"strings"
	"sync"
)

type Dispatcher struct {
	// unexported fields below
	current *item
	queue   *list.List
	mu      sync.RWMutex
}

type item struct {
	cmd  ICommand
	cb   HandleCommandResp
	seen int
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		queue: list.New(),
	}
}

func (s *Dispatcher) Next() ICommand {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Only return the next command if we're not processing one
	if s.current == nil && s.current != nil {
		return nil
	}

	if elem := s.queue.Front(); elem != nil {
		cmd := elem.Value.(*item)
		s.current = cmd
		s.queue.Remove(elem)
		return cmd.cmd
	}

	return nil
}

func (s *Dispatcher) CommandDone() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.current = nil
}

func (s *Dispatcher) OnMsg(msg string) {
	// Strip the trailing newline
	msg = strings.Trim(msg, "")

	s.mu.RLock()
	defer s.mu.RUnlock()

	s.current.seen++
	if s.current == nil || s.current.cmd.SkipFirstMsg() && s.current.seen == 1 {
		return
	}

	s.current.cb(s.current.cmd, msg)
}

func (s *Dispatcher) Len() int { return s.queue.Len() }

func (s *Dispatcher) Enqueue(cmd ICommand, cmdcb HandleCommandResp) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.queue.PushBack(&item{
		cmd: cmd,
		cb:  cmdcb,
	})
}
