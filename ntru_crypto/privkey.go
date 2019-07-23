package ntru_crypto

import (
	"fmt"
	"github.com/azd1997/golang-ntru/ntru_utils/bitpack"
	"github.com/azd1997/golang-ntru/ntru_utils/mgftp1"
	"github.com/azd1997/golang-ntru/ntru_utils/params"
	"github.com/azd1997/golang-ntru/ntru_utils/poly"
)

// A PrivateKey represents a NTRUEncrypt private key.
type PrivateKey struct {
	PublicKey
	F *poly.Polynomial
}

// Size returns the length of the binary representation of this private key.
func (priv *PrivateKey) Size() int {
	commonSize := 1 + len(priv.Params.OIDBytes) + bitpack.PackedLength(int(priv.Params.N), int(priv.Params.Q))
	packedSize := priv.packedSize()
	listedSize := priv.listedSize()
	if priv.packedSize() < priv.listedSize() {
		return commonSize + packedSize
	}
	return commonSize + listedSize
}

// Bytes returns the binary representation of a private key.
func (priv *PrivateKey) Bytes() []byte {
	ret := make([]byte, priv.Size())
	ret[0] = blobPrivateKeyDefaultV1
	copy(ret[1:4], priv.Params.OIDBytes)
	fOff := blobHeaderLen
	fOff += bitpack.Pack(int(priv.Params.N), int(priv.Params.Q), priv.H.Coeffs, 0, ret, fOff)

	F := priv.recoverF()
	if priv.packedSize() < priv.listedSize() {
		// Convert f to a packed F, and write it out.
		fBuf := &bufByteRdWriter{b: ret, off: fOff}
		mgftp1.EncodeTrinomial(F, fBuf)
	} else {
		// Convert f to a listed f.
		bitpack.PackListedCoefficients(F, int(priv.Params.Df), int(priv.Params.Df), ret, fOff)
	}
	F.Reset()

	return ret
}

// NewPrivateKey decodes a PrivateKey from it's binary representation.
func NewPrivateKey(raw []byte) (*PrivateKey, error) {
	priv := &PrivateKey{}

	if len(raw) < blobHeaderLen {
		return nil, fmt.Errorf("ntru: invalid private key blob length")
	}
	if raw[0] != blobPrivateKeyDefaultV1 {
		return nil, fmt.Errorf("ntru: invalid private key blob tag")
	}
	p := params.ParamFromBytes(raw[1:4])
	if p == nil {
		return nil, fmt.Errorf("ntru: unsupported parameter set")
	}
	priv.Params = p

	expLen := 1 + len(priv.Params.OIDBytes) + bitpack.PackedLength(int(priv.Params.N), int(priv.Params.Q))
	packedFLen := int((p.N + 4) / 5)
	packedListedFLen := priv.listedSize()
	if packedFLen < packedListedFLen {
		expLen += packedFLen
	} else {
		expLen += packedListedFLen
	}

	if expLen != len(raw) {
		return nil, fmt.Errorf("ntru: invalid private key blob length")
	}

	// Recover h.
	fOff := blobHeaderLen
	priv.H = poly.NewPoly(int(p.N))
	fOff += bitpack.Unpack(int(p.N), int(p.Q), raw, blobHeaderLen, priv.H.Coeffs, 0)

	// Recover F.
	if packedFLen < packedListedFLen {
		fBuf := &bufByteRdWriter{b: raw, off: fOff}
		priv.F, _ = mgftp1.GenTrinomial(int(p.N), fBuf)
	} else {
		priv.F = poly.NewPoly(int(p.N))
		bitpack.UnpackListedCoefficients(priv.F, int(p.N), int(p.Df), int(p.Df), raw, fOff)
	}

	// Compute f = 1+p*F.
	for i, v := range priv.F.Coeffs {
		priv.F.Coeffs[i] = (p.P * v) & 0xfff
	}
	priv.F.Coeffs[0]++

	return priv, nil
}

// packedSize returns the size of F encoded in the packed format.
func (priv *PrivateKey) packedSize() int {
	return (len(priv.F.Coeffs) + 4) / 5
}

// listedSize returns the size of F encoded in the listed format.
func (priv *PrivateKey) listedSize() int {
	return bitpack.PackedLength(2*int(priv.Params.Df), int(priv.Params.N))
}

// Calculate F = (f - 1) / p.
func (priv *PrivateKey) recoverF() *poly.Polynomial {
	F := poly.NewPoly(len(priv.F.Coeffs))
	F.Coeffs[0] = int16(int8(priv.F.Coeffs[0]-1) / int8(priv.Params.P))
	for i := 1; i < len(F.Coeffs); i++ {
		F.Coeffs[i] = int16(int8(priv.F.Coeffs[i]) / int8(priv.Params.P))
	}
	return F
}

// parseMsgLengthFromM pulls out the message length from a ciphertext.
func (priv *PrivateKey) parseMsgLengthFromM(m []byte) (l int) {
	db := priv.Params.Db >> 3
	if len(m) < int(db+priv.Params.LLen) {
		return -1
	}
	for i := db; i < db+priv.Params.LLen; i++ {
		l = (l << 8) | int(m[i])
	}
	return
}

// verifyMFormat validates that a ciphertext is well formed, and returns the
// message length or -1.
func (priv *PrivateKey) verifyMFormat(m []byte) int {
	ok := true
	db := priv.Params.Db >> 3

	// This is the number of bytes in the formatted message:
	numBytes := db + priv.Params.LLen + int16(priv.Params.MaxMsgLenBytes) + 1
	if len(m) != int(numBytes) {
		ok = false
	}

	// 1) First db bytes are random data.

	// 2) Next lLen bytes are the message length.  Decode and verify.
	//
	// XXX/Yawning This whole block is kind of broken in the Java code.
	//  * It treats the short buffer case as the same as a 0 length msg (ok,
	//    though confusing, since the total length check has failed and ok is
	//    false at this point).
	//  * It checks mLen >= priv.Params.MaxMsgLenBytes, which is blatantly
	//    incorrect and will cause ciphertexts containing maximum length
	//    payload to fail.
	mLen := priv.parseMsgLengthFromM(m)
	if mLen < 0 || mLen > priv.Params.MaxMsgLenBytes {
		// mLen = 1 so that later steps will work (though we will return an
		// error).
		mLen = 1
		ok = false
	}

	// 3) Next mLen bytes are m.

	// 4) Remaining bytes are p0.
	for i := int(db+priv.Params.LLen) + mLen; i < len(m); i++ {
		ok = ok && m[i] == 0
	}

	if ok {
		return mLen
	}
	return -1
}
