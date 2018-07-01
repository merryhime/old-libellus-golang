package objfile

import (
	"compress/zlib"
	"errors"
	"io"
	"strconv"

	"github.com/MerryMage/libellus/objstore/ioutil"
	"github.com/MerryMage/libellus/objstore/objid"
	"github.com/MerryMage/libellus/objstore/objtype"
)

var (
	ErrSizeExceeded = errors.New("objfile: exceeded declared size")
	ErrClosed       = errors.New("objfile: attempted to write to closed writer")
)

type Writer struct {
	z io.WriteCloser
	h ioutil.Hasher
	m io.Writer

	remaining uint64
	closed    bool
}

func NewWriter(inner io.Writer, ot objtype.ObjType, size uint64) (*Writer, error) {
	z := zlib.NewWriter(inner)
	h := ioutil.NewHasher()
	m := io.MultiWriter(z, h)

	// Write Header
	_, err := m.Write([]byte(ot.String()))
	if err != nil {
		return nil, err
	}
	_, err = m.Write([]byte(" "))
	if err != nil {
		return nil, err
	}
	_, err = m.Write([]byte(strconv.FormatUint(size, 10)))
	if err != nil {
		return nil, err
	}
	_, err = m.Write([]byte("\000"))
	if err != nil {
		return nil, err
	}

	return &Writer{
		z:         z,
		h:         h,
		m:         m,
		remaining: size,
		closed:    false,
	}, nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if w.closed {
		return 0, ErrClosed
	}

	if uint64(len(p)) > w.remaining {
		n, _ = w.m.Write(p[0:w.remaining])
		w.remaining = 0
		err = ErrSizeExceeded
		return
	}

	n, err = w.m.Write(p)
	w.remaining -= uint64(n)
	return
}

func (w *Writer) Oid() objid.Oid {
	return w.h.Oid()
}

func (w *Writer) Close() error {
	err := w.z.Close()
	if err != nil {
		return err
	}

	w.closed = true
	return nil
}
