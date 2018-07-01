package commit

import (
	"bytes"
	"io"

	"github.com/MerryMage/libellus/objstore/ioutil"
	"github.com/MerryMage/libellus/objstore/objid"
)

func Read(r io.Reader) (Commit, error) {
	commit := Commit{}

	var eof bool

	for {
		line, err := ioutil.ReadUntil(r, '\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			return commit, err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			// Message comes next
			break
		}

		parts := bytes.SplitN(line, []byte{' '}, 2)
		switch string(parts[0]) {
		case "tree":
			tree, err := objid.FromString(string(parts[1]))
			if err != nil {
				return commit, err
			}
			commit.Tree = tree
		case "parent":
			parent, err := objid.FromString(string(parts[1]))
			if err != nil {
				return commit, err
			}
			commit.Parents = append(commit.Parents, parent)
		case "author":
			author, err := NewSignature(parts[1])
			if err != nil {
				return commit, err
			}
			commit.Author = author
		case "committer":
			committer, err := NewSignature(parts[1])
			if err != nil {
				return commit, err
			}
			commit.Committer = committer
		}

		if eof {
			return commit, nil
		}
	}

	message, err := ioutil.ReadAll(r)
	commit.Message = string(message)
	return commit, err
}
