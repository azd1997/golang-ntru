package ntru_crypto

import "io"


// readerToByteReader 使用一个适配器（adapter）对io.Reader进行包装，实现io.Reader到
// io.ByteReader的转换。或者说允许用户使用io.Reader替代io.ByteReader
// 如果用户自己的代码中已经做了这部分工作，那么这个函数不起作用
func readerToByteReader(r io.Reader) io.ByteReader {
	if br, ok := (r).(io.ByteReader); ok {
		return br
	}
	return &byteReaderAdapter{r: r}
}

type byteReaderAdapter struct {
	r io.Reader
}

// 实现io.ByteReader接口
func (r *byteReaderAdapter) ReadByte() (c byte, err error) {
	// TODO: This is really inefficient, maybe buffer?  Not sure how much I
	//  trust bufio to cleanup, and this is correct so it's ok for now.
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

// bufByteRdWriter is a minimal io.Byte[Read,Writer] implementation that modifies
// an existing slice in place.  Interlacing（交错） ReadByte and WriteByte calls will
// lead to bad things happening.
type bufByteRdWriter struct {
	b   []byte
	off int
}

func (b *bufByteRdWriter) WriteByte(c byte) (err error) {
	if b.off+1 > len(b.b) {
		// This should *NEVER* happen.
		return io.ErrShortWrite
	}
	b.b[b.off] = c
	b.off++
	return nil
}

func (b *bufByteRdWriter) ReadByte() (c byte, err error) {
	if b.off > len(b.b) {
		return 0, io.ErrUnexpectedEOF
	}
	c = b.b[b.off]
	b.off++
	return
}



