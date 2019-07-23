/******************************************************************************
 * NTRU Cryptography Reference Source Code
 * Copyright (c) 2009-2013, by Security Innovation, Inc. All rights reserved.
 *
 * Copyright (C) 2009-2013  Security Innovation
 * Copyright (C) 2014  Yawning Angel (yawning at schwanenlied dot me)
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 *********************************************************************************/

// Package ntru implements the NTRUEncrypt public key cryptosystem.
package ntru_crypto

import (
	"crypto"
	"io"

	"github.com/azd1997/golang-ntru/ntru_utils/bitpack"
	"github.com/azd1997/golang-ntru/ntru_utils/bpgm3"
	"github.com/azd1997/golang-ntru/ntru_utils/igf2"
	"github.com/azd1997/golang-ntru/ntru_utils/params"
	"github.com/azd1997/golang-ntru/ntru_utils/poly"
)

const (
	blobHeaderLen           = 4
	blobPublicKeyV1         = 1
	blobPrivateKeyDefaultV1 = 2
)

var inverterMod2048 poly.Inverter


/*根据指定参数集生成NTRU私钥*/
// random为随机源 (可以使用crypto/rand.Reader).
func GenerateKey(random io.Reader, oid params.Oid) (priv *PrivateKey, err error) {
	//选择参数集
	keyParams := params.Param(oid)
	if keyParams == nil {
		return nil, InvalidParamError(oid)
	}
	//将random(io.Reader)转化为prng(io.ByteReader)
	prng := readerToByteReader(random)
	// 创建igf对象实例，用以生成多项式索引
	igf := igf2.NewFromReader(keyParams.N, keyParams.C, prng)

	// 生成可逆的三项式g
	var g *poly.Polynomial
	for isInvertible := false; !isInvertible; {
		// 使用igf来生成随机三项式，1s = Dg+1; -1s = Dg
		if g, err = bpgm3.GenTrinomial(keyParams.N, keyParams.Dg+1, keyParams.Dg, igf); err != nil {
			return nil, err
		}
		// 判断g是否可逆。
		// 根据多项式可逆的理论，若多项式可逆，则其逆必不为0；多项式g,与g不断移位构成的矩阵G，G*ginv = [1,0,...,0]'
		gInv := inverterMod2048.Invert(g)
		isInvertible = gInv != nil
	}

	// 生成三项式 F, f=1+p*F, and F^-1 mod q.
	var F, f, fInv *poly.Polynomial
	for isInvertible := false; !isInvertible; {
		// 1s = -1s = Df
		if F, err = bpgm3.GenTrinomial(keyParams.N, keyParams.Df, keyParams.Df, igf); err != nil {
			return nil, err
		}
		f = poly.NewPoly(int(keyParams.N))
		for i := range f.Coeffs {
			f.Coeffs[i] = (keyParams.P * F.Coeffs[i]) & 0xfff	//与0xfff按位与，因为int16最大为0xffff，
																// 所以这样做可以使系数最高四位置0
																// TODO：这是为什么？.
		}
		f.Coeffs[0]++	//f[0]+1

		fInv = inverterMod2048.Invert(f)
		isInvertible = fInv != nil
	}

	// 计算 h = f^-1 * g * p mod q.
	h := poly.Convolution(fInv, g)
	for i := range h.Coeffs {
		h.Coeffs[i] = (h.Coeffs[i] * keyParams.P) % keyParams.Q
		if h.Coeffs[i] < 0 {
			h.Coeffs[i] += keyParams.Q
		}
	}

	//fInv、F置零 TODO: 置零的作用是？（猜测：清除内存中这二者信息，提防攻击程序找到这两个求出私钥）
	fInv.Reset()
	F.Reset()

	// 组装公私钥对象
	/*
	type PrivateKey struct {
		PublicKey
		F *poly.Polynomial
	}
	type PublicKey struct {
		Params *params.KeyParams
		H      *poly.Polynomial
	}*/
	priv = &PrivateKey{}
	priv.Params = keyParams
	priv.H = h	// NTRU公钥
	priv.F = f	// NTRU私钥
	return
}

// 对明文进行加密（明文长度不得超过最大明文长度，取决于参数集）
func Encrypt(random io.Reader, pub *PublicKey, msg []byte) (out []byte, err error) {
	// 检查公钥是否正确（包含参数集和H信息）
	if pub.Params == nil || pub.H == nil {
		return nil, ErrInvalidKey
	}
	//检查明文长度是否符合要求
	if pub.Params.MaxMsgLenBytes < len(msg) {
		return nil, ErrMessageTooLong
	}

	// TODO: Get rid of these casts（摆脱这些类型强制转换）.
	var mPrime, R *poly.Polynomial
	for {
		// 组装明文消息 M = b | len | message | p0.
		var M []byte
		M, err = pub.generateM(msg, random)
		if err != nil {
			return nil, err
		}

		// 组装明文三项式 Mtrin = trinary poly derived from M.
		mCoeffs := convPolyBinaryToTrinary(int(pub.Params.N), M)
		Mtrin := poly.NewFromCoeffs(mCoeffs)

		// 组装 sData = OID | m | b | h.
		// 注意！ 组装sData时，m = msg; b = M
		sData := pub.formSData(msg, 0, len(msg), M, 0)

		// 根据 sData 生成三项式 r .
		var r *poly.Polynomial	//根据sData作为种子生成随机索引		TODO:这么大费周折的意义？
		igf := igf2.New(pub.Params.N, pub.Params.C, pub.Params.IGFHash, int(pub.Params.MinCallsR), sData, 0, len(sData))
		r, err = bpgm3.GenTrinomial(pub.Params.N, pub.Params.Dr, pub.Params.Dr, igf)
		if err != nil {
			return nil, err
		}

		// 计算 R = r * h mod q.
		R = poly.ConvolutionModN(r, pub.H, int(pub.Params.Q))

		// 计算 R4 = R mod 4, form octet（八位字节） string.
		// 计算 mask = MGF1(R4, N, minCallsMask).
		mask := pub.calcEncryptionMask(R)

		// 计算 m' = M + mask （mod p）.
		mPrime = Mtrin.AddAndRecenter(mask, int(pub.Params.P), -1)

		// 计数 #1s, #0s, #-1s in m', 若 <dm0 则丢弃.
		if poly.CheckDm0(mPrime, pub.Params.Dm0) {
			break
		}
	}

	// e = R + m' （mod q）.
	e := R.Add(mPrime, int(pub.Params.Q))

	// 将e进行位组合，包装成[]byte
	cLen := bitpack.PackedLength(len(e.Coeffs), int(pub.Params.Q))
	out = make([]byte, cLen)
	bitpack.Pack(len(e.Coeffs), int(pub.Params.Q), e.Coeffs, 0, out, 0)
	return
}

// 解密密文，得到明文
func Decrypt(priv *PrivateKey, ciphertext []byte) (out []byte, err error) {

	// 检查私钥的有效性（是否写入了诸项内容）
	if priv.Params == nil || priv.H == nil || priv.F == nil {
		return nil, ErrInvalidKey
	}

	// TODO: 想办法去除类型的强制转换.
	// 检查密文长度是否符合要求
	expectedCTLength := bitpack.PackedLength(int(priv.Params.N), int(priv.Params.Q))
	if len(ciphertext) != expectedCTLength {
		return nil, ErrDecryption
	}

	fail := false

	// 从密文数据包[]byte中解析出密文三项式[]int16，并再次检查密文长度是否符合要求
	e := poly.NewPoly(int(priv.Params.N))
	numUnpacked := bitpack.Unpack(int(priv.Params.N), int(priv.Params.Q), ciphertext, 0, e.Coeffs, 0)
	if numUnpacked != len(ciphertext) {
		return nil, ErrDecryption
	}

	// a = f * e (mod q) 并将系数限制在[-q/2, q/2)间
	// range [A..A+q-1], where A = lower bound decryption coefficient (-q/2 in all param sets).
	ci := poly.ConvolutionModN(priv.F, e, int(priv.Params.Q))
	for i := range ci.Coeffs {
		if ci.Coeffs[i] >= priv.Params.Q/2 {
			ci.Coeffs[i] -= priv.Params.Q
		}
	}

	// 计算 ci = message candidate = a mod p in [-1, 0, 1].
	for i := 0; i < int(priv.Params.N); i++ {
		ci.Coeffs[i] = int16(int8((ci.Coeffs[i] % priv.Params.P) & 0xff))
		switch ci.Coeffs[i] {
		case 2:
			ci.Coeffs[i] = -1
		case -2:
			ci.Coeffs[i] = 1
		}
	}

	// 检查ci中系数的个数是否符合要求 #1s, #0s < dm0 则失败
	if !poly.CheckDm0(ci, priv.Params.Dm0) {
		fail = true
	}

	// 计算 r*h 候选: cR = e - ci;
	cR := e.Subtract(ci, int(priv.Params.Q))

	// 计算 cR4 = cR mod 4.
	// 调用给定的MGF（掩码生成函数）生成多项式掩码
	mask := priv.calcEncryptionMask(cR)

	// 计算 cMtrin = (cm' - mask) (mod p). 这时候是三进制，或者说三项式
	cMtrin := ci.SubtractAndRecenter(mask, int(priv.Params.P), -1)

	// 将 cMtrin 转换为 cMbin.并丢弃后位.
	cM := priv.convPolyTrinaryToBinary(cMtrin)

	// 解析 cMbin，其组装格式为 b || l || m || p0. 如果格式不匹配，则报错
	mOffset := int(priv.Params.Db/8 + priv.Params.LLen)
	mLen := priv.verifyMFormat(cM)
	if mLen < 0 {
		mLen = 1
		fail = true
	}

	// 由 OID, m, b, hTrunc 组装 sData
	sData := priv.formSData(cM, mOffset, mLen, cM, 0)

	// 从 sData 计算 cr.
	igf := igf2.New(priv.Params.N, priv.Params.C, priv.Params.IGFHash, int(priv.Params.MinCallsR), sData, 0, len(sData))
	cr, err := bpgm3.GenTrinomial(priv.Params.N, priv.Params.Dr, priv.Params.Dr, igf)
	if err != nil {
		fail = true
	}
	igf.Close()

	// 计算 cR' = h * cr mod q
	cRPrime := poly.ConvolutionModN(cr, priv.H, int(priv.Params.Q))
	// If cR != cR', fail
	if !cR.Equals(cRPrime) {
		fail = true
	}

	if fail {
		return nil, ErrDecryption
	}

	out = cM[mOffset : mOffset+mLen]
	return
}


func init() {
	// Q = 2048, prime = 2
	invMod2 := []int16{0, 1}
	inverterMod2048 = poly.NewInverterModPowerOfPrime(2048, 2, invMod2)
}

var _ crypto.PublicKey = (*PublicKey)(nil)
var _ crypto.PrivateKey = (*PrivateKey)(nil)