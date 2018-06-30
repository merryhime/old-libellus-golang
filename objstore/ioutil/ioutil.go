package ioutil

import (
	"io"
)

func ReadUntil(r io.Reader, delim byte) ([]byte, error) {
	var tmp [1]byte
	ret := make([]byte, 0, 32)
	for {
		n, err := r.Read(tmp[:])
		if n > 0 {
			if tmp[0] == delim {
				return ret, nil
			}
			ret = append(ret, tmp[0])
		}
		if err != nil {
			return ret, err
		}
	}
}
