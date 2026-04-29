package backend

import "sync"

type pendingBreakpoint struct {
	hit      BreakpointHit
	respChan chan *MockResponse
}

type BreakpointManager struct {
	mu      sync.Mutex
	pending map[string]*pendingBreakpoint
}

func NewBreakpointManager() *BreakpointManager {
	return &BreakpointManager{pending: make(map[string]*pendingBreakpoint)}
}

// Hold blocks until the breakpoint is released or discarded.
// Returns (true, response) on release, (false, nil) on discard.
func (bm *BreakpointManager) Hold(hit BreakpointHit, onHold func(BreakpointHit)) (bool, *MockResponse) {
	ch := make(chan *MockResponse, 1)
	bm.mu.Lock()
	bm.pending[hit.ID] = &pendingBreakpoint{hit: hit, respChan: ch}
	bm.mu.Unlock()

	onHold(hit)

	resp := <-ch
	return resp != nil, resp
}

func (bm *BreakpointManager) GetPending() []BreakpointHit {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	result := make([]BreakpointHit, 0, len(bm.pending))
	for _, bp := range bm.pending {
		result = append(result, bp.hit)
	}
	return result
}

func (bm *BreakpointManager) Release(id string, resp MockResponse) bool {
	bm.mu.Lock()
	bp, ok := bm.pending[id]
	if ok {
		delete(bm.pending, id)
	}
	bm.mu.Unlock()
	if !ok {
		return false
	}
	bp.respChan <- &resp
	return true
}

func (bm *BreakpointManager) Discard(id string) bool {
	bm.mu.Lock()
	bp, ok := bm.pending[id]
	if ok {
		delete(bm.pending, id)
	}
	bm.mu.Unlock()
	if !ok {
		return false
	}
	bp.respChan <- nil
	return true
}
