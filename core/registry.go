package core

import (
	"fmt"
	"sync"
)

type IRegistry interface {
	Push(Identity, *Pipe)
	Get(Identity) (*Pipe, error)
	Pop(Identity) (*Pipe, error)
	Exists(Identity) bool
	Foreach(callback func(Identity, *Pipe))
}

type Registry struct {
	clients map[Identity]*Pipe
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		clients: make(map[Identity]*Pipe, 0),
	}
}

func (r *Registry) Push(id Identity, pipe *Pipe) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clients[id] = pipe
}

func (r *Registry) Get(id Identity) (*Pipe, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, set := r.clients[id]; !set {
		return nil, fmt.Errorf("client not known: %v", id)
	}

	return r.clients[id], nil
}

func (r *Registry) Pop(id Identity) (*Pipe, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, set := r.clients[id]; !set {
		return nil, fmt.Errorf("client not known: %v", id)
	}

	pipe := r.clients[id]
	delete(r.clients, id)
	return pipe, nil
}

func (r *Registry) Exists(id Identity) bool {
	_, set := r.clients[id]
	return set
}

func (r *Registry) Foreach(callback func(Identity, *Pipe)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for k, v := range r.clients {
		callback(k, v)
	}
}

// TODO Add(Identity, Pipe) bool OR error ?
// TODO Remove(Identity) bool
// TODO Get(Identity) Identity, bool
