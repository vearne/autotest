package model

import (
	"github.com/vearne/autotest/internal/config"
	"sync"
)

type State int

type Reason int

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

func NewStateGroup(testCases []*config.TestCase) *StateGroup {
	var g StateGroup
	g.states = make(map[uint64]State, 0)
	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]
		g.states[tc.ID] = StateNotExecuted
	}
	return &g
}
