package model

import (
	"github.com/fullstorydev/grpcurl"
	"sync"
)

type State int

func (s State) String() string {
	switch s {
	case StateNotExecuted:
		return "StateNotExecuted"
	case StateSuccessFul:
		return "StateSuccessFul"
	case StateFailed:
		return "StateFailed"
	}
	return ""
}

type Reason int

func (r Reason) String() string {
	switch r {
	case ReasonSuccess:
		return "ReasonSuccess"
	case ReasonRequestFailed:
		return "ReasonRequestFailed"
	case ReasonRuleVerifyFailed:
		return "ReasonRuleVerifyFailed"
	case ReasonDependentItemNotCompleted:
		return "ReasonDependentItemNotCompleted"
	case ReasonTemplateRenderError:
		return "ReasonTemplateRenderError"
	case ReasonDependentItemFailed:
		return "ReasonDependentItemFailed"
	}
	return ""
}

type ExportType int

type StateGroup struct {
	states map[uint64]State
	locker sync.RWMutex
}

func (g *StateGroup) SetState(id uint64, s State) {
	g.locker.Lock()
	defer g.locker.Unlock()

	g.states[id] = s
}

func (g *StateGroup) GetState(id uint64) State {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.states[id]
}

type IdItem interface {
	GetID() uint64
}

func NewStateGroup() *StateGroup {
	var g StateGroup
	g.states = make(map[uint64]State, 0)
	return &g
}

type DescSourceCache struct {
	sources map[string]grpcurl.DescriptorSource
	locker  sync.RWMutex
}

func NewDescSourceCache() *DescSourceCache {
	var cache DescSourceCache
	cache.sources = make(map[string]grpcurl.DescriptorSource)
	return &cache

}
func (c *DescSourceCache) Set(target string, s grpcurl.DescriptorSource) {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.sources[target] = s
}

func (c *DescSourceCache) Get(target string) (s grpcurl.DescriptorSource, ok bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	s, ok = c.sources[target]
	return
}

type GrpcResp struct {
	Code    string
	Message string
	Headers []string
	Body    string
}
