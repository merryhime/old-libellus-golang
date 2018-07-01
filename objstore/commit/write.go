package commit

import (
	"fmt"
	"io"
)

func (c Commit) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "tree %s\n", c.Tree)
	if err != nil {
		return err
	}

	for _, parent := range c.Parents {
		_, err := fmt.Fprintf(w, "parent %s\n", parent)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(w, "author %s\n", c.Author)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "committer %s\n", c.Committer)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "\n%s", c.Message)
	if err != nil {
		return err
	}

	return nil
}
