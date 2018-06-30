package objtype

import (
	"fmt"
)

type ObjTypeParseError string

func (e ObjTypeParseError) Error() string {
	return fmt.Sprintf("objtype: Could not parse %q", e)
}

type ObjType int

const (
	Invalid ObjType = iota
	Commit
	Tree
	Blob
	Tag
)

func (ot ObjType) String() string {
	switch ot {
	case Commit:
		return "commit"
	case Tree:
		return "tree"
	case Blob:
		return "blob"
	case Tag:
		return "tag"
	}
	return "invalid"
}

func (ot ObjType) Valid() bool {
	return ot >= Commit && ot <= Tag
}

func Make(s string) (ObjType, error) {
	switch s {
	case "commit":
		return Commit, nil
	case "tree":
		return Tree, nil
	case "blob":
		return Blob, nil
	case "tag":
		return Tag, nil
	}
	return Invalid, ObjTypeParseError(s)
}
