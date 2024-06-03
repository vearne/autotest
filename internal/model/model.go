package model

import (
	"github.com/vearne/autotest/internal/config"
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

func NewStateGroup(testCases []*config.TestCaseHttp) *StateGroup {
	var g StateGroup
	g.states = make(map[uint64]State, 0)
	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]
		g.states[tc.ID] = StateNotExecuted
	}
	return &g
}
