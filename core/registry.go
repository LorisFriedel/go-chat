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
	clients *sync.Map
}

func NewRegistry() *Registry {
	return &Registry{
		clients: &sync.Map{},
	}
}

func (r *Registry) Push(id Identity, pipe *Pipe) {
	r.clients.Store(id, pipe)
}

func (r *Registry) Get(id Identity) (*Pipe, error) {

	p, ok := r.clients.Load(id)
	if !ok {
		return nil, fmt.Errorf("client not known: %v", id)
	}

	return p.(*Pipe), nil
}

func (r *Registry) Pop(id Identity) (*Pipe, error) {
	p, ok := r.clients.Load(id)
	if !ok {
		return nil, fmt.Errorf("client not known: %v", id)
	}

	r.clients.Delete(id)
	return p.(*Pipe), nil
}

func (r *Registry) Exists(id Identity) bool {
	_, ok := r.clients.Load(id)
	return ok
}

func (r *Registry) Foreach(callback func(Identity, *Pipe)) {
	r.clients.Range(func(key, value interface{}) bool {
		callback(key.(Identity), value.(*Pipe))
		return true
	})
}

// TODO Add(Identity, Pipe) bool OR error ?
// TODO Remove(Identity) bool
// TODO Get(Identity) Identity, bool
