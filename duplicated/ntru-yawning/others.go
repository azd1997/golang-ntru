package ntru_yawning

import "io"

// readerToByteReader allows using any io.Reader as a io.ByteReader by wrapping
// the io.Reader with an adapter.  If the reader already implements
// io.ByteReader, then this function is a no-op.
func readerToByteReader(r io.Reader) io.ByteReader {
	if br, ok := (r).(io.ByteReader); ok {
		return br
	}
	return &byteReaderAdapter{r: r}
}

type byteReaderAdapter struct {
	r io.Reader
}

func (r *byteReaderAdapter) ReadByte() (c byte, err error) {
	// TODO: This is really inefficient, maybe buffer?  Not sure how much I
	// trust bufio to cleanup, and this is correct so it's ok for now.
	var b [1]byte
	defer func() {
		b[0] = 0
	}()
	if _, err = r.r.Read(b[:]); err != nil {
		return
	}
	c = b[0]
	return
}