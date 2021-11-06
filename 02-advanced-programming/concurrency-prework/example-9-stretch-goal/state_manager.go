package main

import (
	"fmt"
	"sync"
)

type StateManager struct {
	lock      *sync.RWMutex
	consumers map[int]*Consumer
}

func NewStateManager(numConsumers int) *StateManager {
	s := &StateManager{
		lock:      &sync.RWMutex{},
		consumers: make(map[int]*Consumer),
	}
	for i := 0; i < numConsumers; i++ {
		s.consumers[i] = NewConsumer(i, s)
	}
	return s
}

func (s *StateManager) RemoveConsumer(id int) {
	fmt.Println("removing consumer ", id)
	// s.lock.Lock()
	// defer s.lock.Unlock()
	// hello

	delete(s.consumers, id)
	fmt.Println("removed consumer ", id)
}

func (s *StateManager) GetConsumer(id int) *Consumer {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.consumers[id]
}

func (s *StateManager) PrintState() {
	s.lock.RLock()
	defer s.lock.RUnlock()

	fmt.Println("Started PrintState")
	for _, consumer := range s.consumers {
		fmt.Println(consumer.GetState())
	}
	fmt.Println("Done with PrintState")
}
