package objfile

import (
	"compress/zlib"
	"io"
	"strconv"

	"github.com/MerryMage/libellus/objstore/ioutil"
	"github.com/MerryMage/libellus/objstore/objtype"
)

type Reader struct {
	z    io.ReadCloser
	size uint64
	ot   objtype.ObjType
}

func NewReader(inner io.Reader) (*Reader, error) {
	var z io.ReadCloser
	var raw []byte
	var err error

	z, err = zlib.NewReader(inner)
	if err != nil {
		return nil, err
	}

	// Read Header

	raw, err = ioutil.ReadUntil(z, ' ')
	if err != nil {
		return nil, err
	}

	ot, err := objtype.Make(string(raw))
	if err != nil {
		return nil, err
	}

	raw, err = ioutil.ReadUntil(z, 0)
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseUint(string(raw), 10, 64)
	if err != nil {
		return nil, err
	}

	return &Reader{
		z:    z,
		size: size,
		ot:   ot,
	}, nil
}

func (r *Reader) Size() uint64 {
	return r.size
}

func (r *Reader) ObjType() objtype.ObjType {
	return r.ot
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.z.Read(p)
}

func (r *Reader) Close() error {
	return r.z.Close()
}
