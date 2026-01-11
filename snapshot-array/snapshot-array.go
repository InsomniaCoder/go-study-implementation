package main

import "sort"

type Snapshottable interface {
	Snap() int
	Set(index int, value int)
	Get(snapshotId int, index int) int
}

type Entry struct {
	value      int
	snapshotId int
}

type SnapshotArray struct {
	current    map[int]int
	snapshotId int
	snapshots  map[int][]Entry //store the value at each snapshot
	size       int
}

func NewSnapshotArray(size int) Snapshottable {
	return &SnapshotArray{
		current:    make(map[int]int),
		snapshotId: 0,
		size:       size,
		snapshots:  make(map[int][]Entry),
	}
}

// dirty set
func (sa *SnapshotArray) Set(index int, value int) {
	// on each set, find the snapshot and add it the the history
	sa.current[index] = value
	snapshot := sa.snapshots[index]

	if len(snapshot) > 0 && snapshot[len(snapshot)-1].snapshotId == sa.snapshotId {
		sa.snapshots[index][len(snapshot)-1].value = value
	} else {
		newEntry := Entry{
			snapshotId: sa.snapshotId,
			value:      value,
		}
		sa.snapshots[index] = append(sa.snapshots[index], newEntry)
	}
}

func (sa *SnapshotArray) Snap() int {
	sa.snapshotId++
	return sa.snapshotId - 1
}

func (sa *SnapshotArray) Get(snapshotId int, index int) int {
	//if asking for the latest snapshot
	if snapshotId == sa.snapshotId {
		return sa.current[index]
	}
	//get the history
	snapshot := sa.snapshots[index]
	largestIndex := sort.Search(len(snapshot), func(i int) bool {
		return snapshot[i].snapshotId > snapshotId //return the largest index
	})

	// if it returns n, it means it does not find.
	if largestIndex == 0 {
		return -1
	}

	return snapshot[largestIndex-1].value
}

func main() {

}
