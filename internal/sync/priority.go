package sync

import "sort"

// PriorityLevel represents the execution priority of a profile.
type PriorityLevel int

const (
	PriorityLow    PriorityLevel = 0
	PriorityNormal PriorityLevel = 50
	PriorityHigh   PriorityLevel = 100
)

// PriorityEntry pairs a profile name with its priority level.
type PriorityEntry struct {
	Profile  string
	Priority PriorityLevel
}

// PriorityQueue orders profiles by descending priority for execution.
type PriorityQueue struct {
	entries []PriorityEntry
}

// NewPriorityQueue creates an empty PriorityQueue.
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{}
}

// Add inserts a profile with the given priority into the queue.
func (q *PriorityQueue) Add(profile string, priority PriorityLevel) {
	q.entries = append(q.entries, PriorityEntry{Profile: profile, Priority: priority})
}

// Sorted returns profile names ordered from highest to lowest priority.
// Profiles with equal priority retain insertion order.
func (q *PriorityQueue) Sorted() []string {
	copy := make([]PriorityEntry, len(q.entries))
	for i, e := range q.entries {
		copy[i] = e
	}
	sort.SliceStable(copy, func(i, j int) bool {
		return copy[i].Priority > copy[j].Priority
	})
	names := make([]string, len(copy))
	for i, e := range copy {
		names[i] = e.Profile
	}
	return names
}

// Len returns the number of entries in the queue.
func (q *PriorityQueue) Len() int {
	return len(q.entries)
}

// Clear removes all entries from the queue.
func (q *PriorityQueue) Clear() {
	q.entries = q.entries[:0]
}
