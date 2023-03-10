package rendezvous

import (
	"bytes"
	"crypto/sha256"
	"hash"
	"sort"
	"sync"
)

type Hasher struct {
	mu      *sync.Mutex
	hash    hash.Hash
	order   SortOrder
	members map[string]struct{}
}

// AddMembers adds the given members to the member list
func (h *Hasher) AddMembers(members ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, member := range members {
		h.members[member] = struct{}{}
	}
}

// RemoveMembers removes the given members from the member list
func (h *Hasher) RemoveMembers(members ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, member := range members {
		delete(h.members, member)
	}
}

// SetMembers replaces the current member list with the members provided
func (h *Hasher) SetMembers(members ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.members = make(map[string]struct{}, len(members))
	for _, member := range members {
		h.members[member] = struct{}{}
	}
}

// Member returns the current member list
func (h *Hasher) Members() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	members := make([]string, 0, len(h.members))
	for member := range h.members {
		members = append(members, member)
	}

	sort.Strings(members)
	return members
}

// PriorityList returns all the placement members sorted by their rendezvous hash value for the given key.
func (h *Hasher) Prioritise(key string) []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.priorityList(key)
}

// Place returns the top n members for the given key sorted by their rendezvous hash value.
func (h *Hasher) Place(key string, n int) []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.place(key, n)
}

// Owner returns the best member for the given key sorted by their rendezvous hash value.
func (h *Hasher) Owner(key string) string {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.priorityList(key)[0]
}

// priorityList returns all the placement members sorted by their rendezvous hash value for the given key.
func (h *Hasher) priorityList(key string) []string {
	// If there are no members, return an empty list
	if len(h.members) == 0 {
		return []string{}
	}

	// If there is only one member, return a slice with that member
	if len(h.members) == 1 {
		for member := range h.members {
			return []string{member}
		}
	}

	// build a list of members
	members := make([]string, 0, len(h.members))
	for member := range h.members {
		members = append(members, member)
	}

	// Use the sorting package to sort the nodes by their rendezvous hash value
	sort.Slice(members, func(i, j int) bool {
		h.hash.Reset()
		h.hash.Write([]byte(members[i]))
		h.hash.Write([]byte(key))
		iHash := h.hash.Sum(nil)

		h.hash.Reset()
		h.hash.Write([]byte(members[j]))
		h.hash.Write([]byte(key))
		jHash := h.hash.Sum(nil)

		return bytes.Compare(iHash, jHash) < 0
	})

	// If we need highest first, reverse the slice
	if h.order == HighestFirst {
		for i := len(members)/2 - 1; i >= 0; i-- {
			opp := len(members) - 1 - i
			members[i], members[opp] = members[opp], members[i]
		}
	}

	// Return the sorted members
	return members
}

func (h *Hasher) place(key string, n int) []string {
	// If there are no members, or n is less than or equal to zero, return an empty slice
	if len(h.members) == 0 || n <= 0 {
		return []string{}
	}

	// Constrain the upper bound of n to the number of members
	if n > len(h.members) {
		n = len(h.members)
	}

	// Get the sorted members and return the first n
	return h.priorityList(key)[:n]
}

func NewHasher(opts ...HasherOption) *Hasher {
	rtn := &Hasher{
		mu:      &sync.Mutex{},
		hash:    sha256.New(), // use sha 256 as the default hash function
		order:   HighestFirst, // default to highest first as it is more common
		members: map[string]struct{}{},
	}

	for _, opt := range opts {
		opt(rtn)
	}
	return rtn
}

type HasherOption func(*Hasher)

func WithHashImplementation(h hash.Hash) HasherOption {
	return func(hs *Hasher) {
		hs.hash = h
	}
}

func WithSortOrder(order SortOrder) HasherOption {
	return func(hs *Hasher) {
		hs.order = order
	}
}

func WithMembers(members ...string) HasherOption {
	return func(hs *Hasher) {
		for _, member := range members {
			hs.members[member] = struct{}{}
		}
	}
}

type SortOrder bool

const (
	LowestFirst  SortOrder = true
	HighestFirst SortOrder = false
)
