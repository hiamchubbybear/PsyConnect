



package idle

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)


var timeAfterFunc = func(d time.Duration, f func()) *time.Timer {
	return time.AfterFunc(d, f)
}



type Enforcer interface {
	ExitIdleMode() error
	EnterIdleMode()
}



type Manager struct {
	
	lastCallEndTime           int64 
	activeCallsCount          int32 
	activeSinceLastTimerCheck int32 
	closed                    int32 

	
	
	enforcer Enforcer 
	timeout  time.Duration

	
	
	
	
	
	
	
	
	
	
	
	idleMu       sync.RWMutex
	actuallyIdle bool
	timer        *time.Timer
}



func NewManager(enforcer Enforcer, timeout time.Duration) *Manager {
	return &Manager{
		enforcer:         enforcer,
		timeout:          timeout,
		actuallyIdle:     true,
		activeCallsCount: -math.MaxInt32,
	}
}



func (m *Manager) resetIdleTimerLocked(d time.Duration) {
	if m.isClosed() || m.timeout == 0 || m.actuallyIdle {
		return
	}

	
	
	if m.timer != nil {
		m.timer.Stop()
	}
	m.timer = timeAfterFunc(d, m.handleIdleTimeout)
}

func (m *Manager) resetIdleTimer(d time.Duration) {
	m.idleMu.Lock()
	defer m.idleMu.Unlock()
	m.resetIdleTimerLocked(d)
}




func (m *Manager) handleIdleTimeout() {
	if m.isClosed() {
		return
	}

	if atomic.LoadInt32(&m.activeCallsCount) > 0 {
		m.resetIdleTimer(m.timeout)
		return
	}

	
	
	if atomic.LoadInt32(&m.activeSinceLastTimerCheck) == 1 {
		
		
		atomic.StoreInt32(&m.activeSinceLastTimerCheck, 0)
		m.resetIdleTimer(time.Duration(atomic.LoadInt64(&m.lastCallEndTime)-time.Now().UnixNano()) + m.timeout)
		return
	}

	
	
	if m.tryEnterIdleMode() {
		
		return
	}

	
	
	
	m.resetIdleTimer(m.timeout)
}








func (m *Manager) tryEnterIdleMode() bool {
	
	
	if !atomic.CompareAndSwapInt32(&m.activeCallsCount, 0, -math.MaxInt32) {
		
		
		
		
		return false
	}
	
	

	m.idleMu.Lock()
	defer m.idleMu.Unlock()

	if atomic.LoadInt32(&m.activeCallsCount) != -math.MaxInt32 {
		
		atomic.AddInt32(&m.activeCallsCount, math.MaxInt32)
		return false
	}
	if atomic.LoadInt32(&m.activeSinceLastTimerCheck) == 1 {
		
		
		
		atomic.AddInt32(&m.activeCallsCount, math.MaxInt32)
		return false
	}

	
	
	
	m.enforcer.EnterIdleMode()
	m.actuallyIdle = true
	return true
}


func (m *Manager) EnterIdleModeForTesting() {
	m.tryEnterIdleMode()
}


func (m *Manager) OnCallBegin() error {
	if m.isClosed() {
		return nil
	}

	if atomic.AddInt32(&m.activeCallsCount, 1) > 0 {
		
		atomic.StoreInt32(&m.activeSinceLastTimerCheck, 1)
		return nil
	}

	
	
	if err := m.ExitIdleMode(); err != nil {
		
		
		atomic.AddInt32(&m.activeCallsCount, -1)
		return err
	}

	atomic.StoreInt32(&m.activeSinceLastTimerCheck, 1)
	return nil
}



func (m *Manager) ExitIdleMode() error {
	
	m.idleMu.Lock()
	defer m.idleMu.Unlock()

	if m.isClosed() || !m.actuallyIdle {
		
		
		
		
		
		
		
		
		
		
		
		return nil
	}

	if err := m.enforcer.ExitIdleMode(); err != nil {
		return fmt.Errorf("failed to exit idle mode: %w", err)
	}

	
	atomic.AddInt32(&m.activeCallsCount, math.MaxInt32)
	m.actuallyIdle = false

	
	m.resetIdleTimerLocked(m.timeout)
	return nil
}


func (m *Manager) OnCallEnd() {
	if m.isClosed() {
		return
	}

	
	atomic.StoreInt64(&m.lastCallEndTime, time.Now().UnixNano())

	
	
	
	
	atomic.AddInt32(&m.activeCallsCount, -1)
}

func (m *Manager) isClosed() bool {
	return atomic.LoadInt32(&m.closed) == 1
}


func (m *Manager) Close() {
	atomic.StoreInt32(&m.closed, 1)

	m.idleMu.Lock()
	if m.timer != nil {
		m.timer.Stop()
		m.timer = nil
	}
	m.idleMu.Unlock()
}
