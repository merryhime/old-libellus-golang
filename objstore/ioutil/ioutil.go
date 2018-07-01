package ioutil

import (
	"bytes"
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

func ReadAll(r io.Reader) ([]byte, error) {
	var b bytes.Buffer
	_, err := b.ReadFrom(r)
	return b.Bytes(), err
}
