package main

import (
	"sync"
	"time"
)

// SetEntry represents one rest-then-pause cycle: the intentional rest
// duration that just elapsed (PreSetRest), and the "slack" duration -
// how long you took between pressing Pause and pressing Resume again.
type SetEntry struct {
	Index      int
	PreSetRest time.Duration // frozen duration of the rest segment before this pause
	SlackStart time.Time     // wall-clock moment this pause happened
	Slack      time.Duration // locked value once Resume is pressed
	Locked     bool          // true once Resume has been pressed for this entry
}

// AppState holds all timer state. The main timer never resets on its own;
// it only accumulates time while "running" (resting) and freezes while
// paused (lifting/slacking). It only goes back to zero via Reset().
type AppState struct {
	mu sync.Mutex

	running      bool
	accumulated  time.Duration // sum of all completed resting segments
	segmentStart time.Time     // when the current running segment began (valid if running)
	started      bool          // true once Start has ever been pressed (since last Reset)

	entries []*SetEntry
}

func NewAppState() *AppState {
	return &AppState{}
}

// TotalElapsed returns the current value of the main "Total Workout Action
// Time" display: accumulated rest time, plus whatever has elapsed in the
// currently-running segment (if any).
func (s *AppState) TotalElapsed(now time.Time) time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	total := s.accumulated
	if s.running {
		total += now.Sub(s.segmentStart)
	}
	return total
}

func (s *AppState) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *AppState) HasStarted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.started
}

// CurrentSetNumber returns the set number to display: while resting, it's
// the set you're about to do (len(entries)+1); while paused (lifting), it's
// the set currently in progress (len(entries)).
func (s *AppState) CurrentSetNumber() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return len(s.entries) + 1
	}
	return len(s.entries)
}

// Start begins (or resumes) resting. If there was an active pause, it locks
// in that pause's slack duration.
func (s *AppState) Start(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return
	}
	if n := len(s.entries); n > 0 {
		last := s.entries[n-1]
		if !last.Locked {
			last.Slack = now.Sub(last.SlackStart)
			last.Locked = true
		}
	}
	s.running = true
	s.started = true
	s.segmentStart = now
}

// Pause ends the current resting segment, folds it into the accumulated
// total, and opens a new (unlocked) log entry whose slack starts ticking now.
func (s *AppState) Pause(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	segDur := now.Sub(s.segmentStart)
	s.accumulated += segDur
	s.running = false

	entry := &SetEntry{
		Index:      len(s.entries) + 1,
		PreSetRest: segDur,
		SlackStart: now,
	}
	s.entries = append(s.entries, entry)
}

// Toggle flips between Start and Pause depending on current state.
func (s *AppState) Toggle(now time.Time) {
	if s.IsRunning() {
		s.Pause(now)
	} else {
		s.Start(now)
	}
}

// Reset wipes everything back to the initial idle state.
func (s *AppState) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
	s.started = false
	s.accumulated = 0
	s.segmentStart = time.Time{}
	s.entries = nil
}

// EntryCount returns how many log entries currently exist.
func (s *AppState) EntryCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}

// EntrySnapshot describes an entry's current display values at time `now`.
type EntrySnapshot struct {
	Index      int
	PreSetRest time.Duration
	Slack      time.Duration
	InProgress bool
}

// Entry returns a snapshot of entry i (0-indexed), with live slack computed
// if it's still unlocked.
func (s *AppState) Entry(i int, now time.Time) EntrySnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.entries[i]
	slack := e.Slack
	if !e.Locked {
		slack = now.Sub(e.SlackStart)
	}
	return EntrySnapshot{
		Index:      e.Index,
		PreSetRest: e.PreSetRest,
		Slack:      slack,
		InProgress: !e.Locked,
	}
}

// Stats returns (averageSlackOfLockedEntries, totalSlackIncludingLive).
func (s *AppState) Stats(now time.Time) (avg, total time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.entries) == 0 {
		return 0, 0
	}
	var lockedSum time.Duration
	var lockedCount int
	var liveSum time.Duration
	for _, e := range s.entries {
		if e.Locked {
			lockedSum += e.Slack
			lockedCount++
		} else {
			liveSum += now.Sub(e.SlackStart)
		}
	}
	total = lockedSum + liveSum
	if lockedCount > 0 {
		avg = lockedSum / time.Duration(lockedCount)
	}
	return avg, total
}
