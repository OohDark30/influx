package execute

import (
        "sort"

        "github.com/EMCECS/influx/query"
)

type GroupLookup struct {
        // groups contains groups of group keys in sorted order.
        // These are optimized for appending access.
        groups []*groupKeyList

        // lastIndex contains the last group that an entry was
        // found in or appended to. This is used to optimize appending.
        lastIndex int

        // nextID is the next id that will be assigned to a key group.
        nextID int
}

// groupKeyList is a group of keys in sorted order.
type groupKeyList struct {
        id       int // unique id for the key group within the group lookup
        elements []groupKeyListElement
        deleted  int
}

type groupKeyListElement struct {
        key     query.GroupKey
        value   interface{}
        deleted bool
}

func (kg *groupKeyList) First() query.GroupKey {
        return kg.elements[0].key
}

func (kg *groupKeyList) Last() query.GroupKey {
        return kg.elements[len(kg.elements)-1].key
}

func (kg *groupKeyList) set(i int, value interface{}) {
        if kg.elements[i].deleted {
                kg.elements[i].deleted = false
                kg.deleted--
        }
        kg.elements[i].value = value
}

func (kg *groupKeyList) delete(i int) {
        kg.elements[i].value = nil
        kg.elements[i].deleted = true
        kg.deleted++
}

// Index determines the location of this key within the key group.
// It returns -1 if this key does not exist within the group.
// It will return -1 if the entry is present, but deleted.
func (kg *groupKeyList) Index(key query.GroupKey) int {
        i := kg.InsertAt(key)
        if i >= len(kg.elements) || kg.elements[i].deleted || !kg.elements[i].key.Equal(key) {
                return -1
        }
        return i
}

// InsertAt will return the index where this key should be inserted.
// If this key would be inserted before the first element, this will
// return 0. If the element exists, then this will return the index
// where that element is located. If the key should be inserted at the
// end of the array, it will return an index the size of the array.
func (kg *groupKeyList) InsertAt(key query.GroupKey) int {
        if kg.Last().Less(key) {
                return len(kg.elements)
        }
        return sort.Search(len(kg.elements), func(i int) bool {
                return !kg.elements[i].key.Less(key)
        })
}

func (kg *groupKeyList) At(i int) interface{} {
        return kg.elements[i].value
}

func NewGroupLookup() *GroupLookup {
        return &GroupLookup{
                lastIndex: -1,
                nextID:    1,
        }
}

func (l *GroupLookup) Lookup(key query.GroupKey) (interface{}, bool) {
        if key == nil || len(l.groups) == 0 {
                return nil, false
        }

        group := l.lookupGroup(key)
        if group == -1 {
                return nil, false
        }

        i := l.groups[group].Index(key)
        if i != -1 {
                return l.groups[group].At(i), true
        }
        return nil, false
}

func (l *GroupLookup) Set(key query.GroupKey, value interface{}) {
        group := l.lookupGroup(key)
        l.createOrSetInGroup(group, key, value)
}

// lookupGroup finds the group index where this key would be located if it were to
// be found or inserted. If no suitable group can be found, then this will return -1
// which indicates that a group has to be created at index 0.
func (l *GroupLookup) lookupGroup(key query.GroupKey) int {
        if l.lastIndex >= 0 {
                kg := l.groups[l.lastIndex]
                if !key.Less(kg.First()) {
                        // If the next group doesn't exist or has a first value that is
                        // greater than this key, then we can return the last index and
                        // avoid performing a binary search.
                        if l.lastIndex == len(l.groups)-1 || key.Less(l.groups[l.lastIndex+1].First()) {
                                return l.lastIndex
                        }
                }
        }

        // Find the last group where the first key is less than or equal
        // than the key we are looking for. This means we need to search for
        // the first group where the first key is greater than the key we are setting
        // and use the group before that one.
        index := sort.Search(len(l.groups), func(i int) bool {
                return key.Less(l.groups[i].First())
        }) - 1
        if index >= 0 {
                l.lastIndex = index
        }
        return index
}

func (l *GroupLookup) createOrSetInGroup(index int, key query.GroupKey, value interface{}) {
        // If this index is at -1, then we are inserting a value with a smaller key
        // than every group and we need to create a new group to insert it at the
        // beginning.
        if index == -1 {
                l.groups = append(l.groups, nil)
                copy(l.groups[1:], l.groups[:])
                l.groups[0] = l.newKeyGroup([]groupKeyListElement{
                        {key: key, value: value},
                })
                l.lastIndex = 0
                return
        }

        kg := l.groups[index]

        // Find the location where this should be inserted.
        i := kg.InsertAt(key)

        // If this should be inserted after the last element, do it and leave.
        if i == len(kg.elements) {
                kg.elements = append(kg.elements, groupKeyListElement{
                        key:   key,
                        value: value,
                })
                return
        } else if kg.elements[i].key.Equal(key) {
                // If the entry already exists at this index, set the value.
                kg.set(i, value)
                return
        }

        // We have to split this entry into two new elements. First, we start
        // by creating space for the new entry.
        l.groups = append(l.groups, nil)
        copy(l.groups[index+2:], l.groups[index+1:])
        // Construct the new group entry and copy the end of the slice
        // into the new key group.
        l.groups[index+1] = func() *groupKeyList {
                entries := make([]groupKeyListElement, len(kg.elements[i:]))
                copy(entries, kg.elements[i:])
                return l.newKeyGroup(entries)
        }()
        // Use a slice on the key group elements to remove the extra elements.
        // Then append the new key group entry.
        kg.elements = kg.elements[:i:cap(kg.elements)]
        kg.elements = append(kg.elements, groupKeyListElement{
                key:   key,
                value: value,
        })
}

func (l *GroupLookup) newKeyGroup(entries []groupKeyListElement) *groupKeyList {
        id := l.nextID
        l.nextID++
        return &groupKeyList{
                id:       id,
                elements: entries,
        }
}

func (l *GroupLookup) Delete(key query.GroupKey) (v interface{}, found bool) {
        if key == nil {
                return
        }

        group := l.lookupGroup(key)
        if group == -1 {
                return nil, false
        }

        kg := l.groups[group]
        i := kg.Index(key)
        if i == -1 {
                return nil, false
        }
        v = kg.At(i)
        kg.delete(i)
        if len(kg.elements) == kg.deleted {
                // All elements in this have been deleted so delete this node.
                copy(l.groups[group:], l.groups[group+1:])
                l.groups = l.groups[: len(l.groups)-1 : cap(l.groups)]
                l.lastIndex = -1
        }
        return v, true
}

// Range will iterate over all groups keys in a stable ordering.
// Range must not be called within another call to Range.
// It is safe to call Set/Delete while ranging.
func (l *GroupLookup) Range(f func(key query.GroupKey, value interface{})) {
        for i := 0; i < len(l.groups); {
                kg := l.groups[i]
                for j := 0; j < len(kg.elements); j++ {
                        entry := kg.elements[j]
                        if entry.deleted {
                                continue
                        }
                        f(entry.key, entry.value)
                }
                if i < len(l.groups) && l.groups[i].id == kg.id {
                        i++
                }
        }
}
