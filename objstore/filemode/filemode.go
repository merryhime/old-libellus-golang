package filemode

import (
	"fmt"
	"strconv"
)

type FileMode int

const (
	Invalid    FileMode = 0
	Dir        FileMode = 0040000
	Regular    FileMode = 0100644
	Executable FileMode = 0100755
	Symlink    FileMode = 0120000
	Submodule  FileMode = 0160000
)

func (fm FileMode) String() string {
	return fmt.Sprintf("%07o", uint32(fm))
}

func (fm FileMode) Valid() bool {
	switch fm {
	case Dir, Regular, Executable, Symlink, Submodule:
		return true
	}
	return false
}

func New(s string) (FileMode, error) {
	n, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return Invalid, err
	}
	return FileMode(n), nil
}
