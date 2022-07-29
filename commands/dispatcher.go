package commands

import (
	"container/list"
	"sync"
)

type Dispatcher struct {
	queue *list.List
	mu    sync.Mutex
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		queue: list.New(),
	}
}

func (s *Dispatcher) Next() ICommand {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item := s.queue.Front(); item != nil {
		return item.Value.(ICommand)
	}

	return nil
}

func (s *Dispatcher) Len() int { return s.queue.Len() }

func (s *Dispatcher) Enqueue(cmd ICommand) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.queue.PushBack(cmd)
}
