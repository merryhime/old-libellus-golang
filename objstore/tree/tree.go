package tree

import (
	"sort"
)

type Tree struct {
	Entries []Entry
}

func (t *Tree) Sort() {
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].sortName() < t.Entries[j].sortName()
	})
}

func (t *Tree) Find(name string) *Entry {
	for i := range t.Entries {
		if name == t.Entries[i].Name {
			return &t.Entries[i]
		}
	}
	return nil
}

func (t *Tree) Add(e Entry) error {
	if t.Find(e.Name) != nil {
		return NameAlreadyExistsError(e.Name)
	}
	t.Entries = append(t.Entries, e)
	return nil
}

func (t *Tree) Delete(name string) error {
	for i := range t.Entries {
		if name == t.Entries[i].Name {
			t.Entries = append(t.Entries[:i], t.Entries[i+1:]...)
			return nil
		}
	}
	return NotFoundError(name)
}
