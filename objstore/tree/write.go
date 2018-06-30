package tree

import (
	"fmt"
	"io"
)

func (t *Tree) Write(w io.Writer) error {
	t.Sort()

	for _, e := range t.Entries {
		_, err := fmt.Fprintf(w, "%o %s\000", e.Mode, e.Name)
		if err != nil {
			return err
		}

		err = e.Oid.Write(w)
		if err != nil {
			return err
		}
	}
	return nil
}
