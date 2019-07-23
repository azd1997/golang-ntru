package ntru_crypto

import (
	"fmt"
	"github.com/azd1997/golang-ntru/ntru_utils/bitpack"
	"github.com/azd1997/golang-ntru/ntru_utils/mgf1"
	"github.com/azd1997/golang-ntru/ntru_utils/mgftp1"
	"github.com/azd1997/golang-ntru/ntru_utils/params"
	"github.com/azd1997/golang-ntru/ntru_utils/poly"
	"io"
)

// A PublicKey represents a NTRUEncrypt public key.
type PublicKey struct {
	Params *params.KeyParams
	H      *poly.Polynomial
}

// Size returns the length of the binary representation of this public key.
func (pub *PublicKey) Size() int {
	return 1 + len(pub.Params.OIDBytes) + bitpack.PackedLength(int(pub.Params.N), int(pub.Params.Q))
}

// Bytes returns the binary representation of a public key.
func (pub *PublicKey) Bytes() []byte {
	ret := make([]byte, pub.Size())
	ret[0] = blobPublicKeyV1
	copy(ret[1:4], pub.Params.OIDBytes)
	bitpack.Pack(int(pub.Params.N), int(pub.Params.Q), pub.H.Coeffs, 0, ret, blobHeaderLen)
	return ret
}

// NewPublicKey decodes a PublicKey from it's binary representation.
func NewPublicKey(raw []byte) (*PublicKey, error) {
	if len(raw) < blobHeaderLen {
		return nil, fmt.Errorf("ntru: invalid public key blob length")
	}
	if raw[0] != blobPublicKeyV1 {
		return nil, fmt.Errorf("ntru: invalid public key blob tag")
	}
	p := params.ParamFromBytes(raw[1:4])
	if p == nil {
		return nil, fmt.Errorf("ntru: unsupported parameter set")
	}

	packedHLen := bitpack.UnpackedLength(int(p.N), int(p.Q))
	if blobHeaderLen+packedHLen != len(raw) {
		return nil, fmt.Errorf("ntru: invalid public key blob length")
	}

	h := poly.NewPoly(int(p.N))
	bitpack.Unpack(int(p.N), int(p.Q), raw, blobHeaderLen, h.Coeffs, 0)
	return &PublicKey{Params: p, H: h}, nil
}

// 根据明文组装完整的明文数据包 M = b | mLen | m | p0.
func (pub *PublicKey) generateM(msg []byte, rng io.Reader) (m []byte, err error) {
	db := pub.Params.Db >> 3	// 右移3位。相当于Db/8，求出存储Db所占长度
	mLen := db + pub.Params.LLen + int16(pub.Params.MaxMsgLenBytes) + 1	//TODO：为何+1
	// +1也就是多出一个字节存储

	m = make([]byte, mLen)
	// 按字节读m
	if _, err = rng.Read(m); err != nil {
		return nil, err
	}

	m[db] = byte(len(msg))		// 存储明文长度信息； m[db]指向的是明文消息长度的存储区域
	copy(m[db+pub.Params.LLen:], msg)	// 存储明文信息
	// 若明文长度未达最大长度，则余者全部补0
	for i := db + pub.Params.LLen + int16(len(msg)); i < mLen; i++ {
		m[i] = 0
	}
	return
}

// convPolyTrinaryToBinary converts a polynomial to a bit-packed binary
// array.
func (pub *PublicKey) convPolyTrinaryToBinary(trin *poly.Polynomial) (b []byte) {
	// The output of this operation is supposed to have the form
	// (b | mLen | m | p0) so we can calculate how many bytes that is supposed
	// to be.

	numBytes := int(pub.Params.Db/8) + int(pub.Params.LLen) + pub.Params.MaxMsgLenBytes + 1
	b = make([]byte, numBytes)
	i, j := 0, 0
	for j < numBytes {
		convPolyTrinaryToBinaryBlockHelper(i, trin.Coeffs, j, b)
		i += 16
		j += 3
	}
	return
}

// formSData 组装字节序列 sData = <OID | m | b | hTrunc>
// hTrunc 为公钥h的位组合表示的前缀（prefix of the bit-packed representation of the public key h）
// mOffset、bOffset分别为明文数据包和b数据包中的实际内容的索引偏移量.
func (pub *PublicKey) formSData(m []byte, mOffset, mLen int, b []byte, bOffset int) (sData []byte) {

	bLen := int(pub.Params.Db >> 3)		//记录区块长度所用字节数
	hLen := int(pub.Params.PkLen >> 3)	//记录公钥长度所用字节数，用以从公钥中提取其前缀（长度信息）

	offset := 0		//索引偏移量
	sDataLen := len(pub.Params.OIDBytes)+mLen+bLen+hLen
	sData = make([]byte, sDataLen)

	// SData装入pub.Params.OIDBytes
	copy(sData[offset:], pub.Params.OIDBytes)
	offset += len(pub.Params.OIDBytes)

	// SData装入明文数据包中的明文消息（三项式表示形式）
	copy(sData[offset:], m[mOffset:mOffset+mLen])
	offset += mLen

	// SData装入b数据包中的b消息
	copy(sData[offset:], b[bOffset:bOffset+bLen])
	offset += bLen

	// 将 pub.H.Coeffs由[]int16转为字节数组[]byte存储,并添加到SData后边
	bitpack.PackN(int(pub.Params.N), int(pub.Params.Q), hLen, pub.H.Coeffs, 0, sData, offset)
	return	//sData = < OID | m | b | h >
}

// calcEncryptionMask calculates the trinomial 'mask' using a bit-packed 'R mod
// 4' as the seed of the MGF_TP_1 algorithm.
func (pub *PublicKey) calcEncryptionMask(r *poly.Polynomial) (p *poly.Polynomial) {
	var err error
	r4 := poly.CalcPolyMod4Packed(r)
	mgf := mgf1.New(pub.Params.MGFHash, int(pub.Params.MinCallsMask), true, r4, 0, len(r4))
	defer mgf.Close()
	p, err = mgftp1.GenTrinomial(int(pub.Params.N), mgf)
	if err != nil {
		panic(err)
	}
	return
}