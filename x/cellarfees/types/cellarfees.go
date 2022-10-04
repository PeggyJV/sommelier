package types

import (
	"sort"
	"strings"
)

// Implementing sort.Interface (see https://pkg.go.dev/sort#Interface)
func (f FeeAccrualCounters) Len() int {
	return len(f.Counters)
}

func (f FeeAccrualCounters) Less(i, j int) bool {
	return strings.Compare(f.Counters[i].Denom, f.Counters[j].Denom) == -1
}

func (f FeeAccrualCounters) Swap(i, j int) {
	f.Counters[i], f.Counters[j] = f.Counters[j], f.Counters[i]
}

// WARNING: If editing these methods, be aware that insert() doesn't check if the denom is already
// present. Duplicate denom entries can result in fees not being auctioned.
func (f *FeeAccrualCounters) insertCounter(denom string, count uint64) {
	f.Counters = append(f.Counters, FeeAccrualCounter{Denom: denom, Count: count})
	sort.Sort(f)
}

// Increment fee accrual counter for denom. If it isn't present, append it to the counters
// slice and increment to 1.
func (f *FeeAccrualCounters) IncrementCounter(denom string) uint64 {
	found := false
	var count uint64
	for i, k := range f.Counters {
		if k.Denom == denom {
			found = true
			f.Counters[i].Count += 1
			count = f.Counters[i].Count
			break
		}
	}
	if !found {
		f.insertCounter(denom, 1)
		count = 1
	}

	return count
}

// Sets the denom's fee accrual counter to zero
func (f *FeeAccrualCounters) ResetCounter(denom string) {
	found := false
	for i, k := range f.Counters {
		if k.Denom == denom {
			found = true
			f.Counters[i].Count = 0
			break
		}
	}
	if !found {
		f.insertCounter(denom, 0)
	}
}
