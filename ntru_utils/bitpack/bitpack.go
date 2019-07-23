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
// TODO:搞懂bitpack在干嘛，是转换还是添加。大概率是转换

// Package bitpack 将 []int16 转换成 []byte
package bitpack

import (
	"github.com/azd1997/golang-ntru/ntru_utils/poly"
)

// CountBits 计算表示一个值所需的最小位数
func CountBits(value int) uint {
	for i := uint(0); i < 32; i++ {
		if (1 << i) > value {
			return i
		}
	}
	return 32
}

// PackedLength 计算位组合输出结果存储所需的字节长度
func PackedLength(numElts, maxEltValue int) int {
	return Pack(numElts, maxEltValue, nil, 0, nil, 0)
}

// Pack 将[]int16 进行位组合，转换成 []byte
func Pack(numElts, maxEltValue int, src []int16, srcOffset int, tgt []byte, tgtOffset int) int {
	bitsPerElement := int(CountBits(maxEltValue - 1))	//TODO:为什么-1
	maxOutLen := (numElts*bitsPerElement + 7) / 8	//确保字节数足够
	return PackN(numElts, maxEltValue, maxOutLen, src, srcOffset, tgt, tgtOffset)
}

// PackN 将[]int16 进行位组合，转换成 []byte并添加到目标[]byte之后，
// 当预定义（N个，maxOutLen）的字节数被生成之后停止
// Elts := Elements
func PackN(numElts, maxEltValue, maxOutLen int, src []int16, srcOffset int, tgt []byte, tgtOffset int) int {

	// 目标字节数组应不为空。（目标数组指待打包的数据）
	if tgt == nil {
		return maxOutLen
	}

	//计算最大元素值所需的最小位数，作为元素的表示位数
	bitsPerElement := CountBits(maxEltValue - 1)

	//取索引上下限
	i, iMax := srcOffset, srcOffset+numElts-1
	j, jMax := tgtOffset, tgtOffset+maxOutLen
	var cur byte
	next := int32(src[i])
	cb, nb := uint(0), bitsPerElement
	for j < jMax && (i < iMax || cb+nb > 8) {
		if cb+nb < 8 {
			// Accumulate next into cur.  The REsult will still be less than 8
			// bits.  Then update next will the next input value.
			cur |= byte(next << (8 - cb - nb))
			cb += nb
			i++
			next = 0x0ffff & int32(src[i])
			nb = bitsPerElement
		} else {
			// Pull the most significant bits off next into cur to make cur 8
			// bits and save it tn the output stream.  Then clear cur, and mask
			// the used bits out of next.
			shift := cb + nb - 8
			// tmp := 0xff & (cur | byte(next >> shift))
			tgt[j] = cur | byte(next>>shift)
			j++
			cur, cb = 0, 0
			next &= lowBitMask(shift)
			nb = shift
		}
	}
	if j < jMax {
		tgt[j] = byte(next << (8 - nb))
		j++
	}
	return j - tgtOffset
}

// PackListedCoefficients bit-packs a polynomial into a listed representation.
func PackListedCoefficients(f *poly.Polynomial, numOnes, numNegOnes int, out []byte, offset int) int {
	if out == nil {
		return PackedLength(numOnes+numNegOnes, len(f.Coeffs))
	}

	coefficients := make([]int16, numOnes+numNegOnes)
	ones, negOnes := 0, numOnes
	for i, v := range f.Coeffs {
		switch v {
		case 1:
			coefficients[ones] = int16(i)
			ones++
		case -1:
			coefficients[negOnes] = int16(i)
			negOnes++
		}
	}
	bpe := int(CountBits(len(f.Coeffs) - 1))
	maxL := ((numOnes+numNegOnes)*bpe + 7) / 8
	l := PackN(numOnes+numNegOnes, len(f.Coeffs), maxL, coefficients, 0, out, offset)

	for i := range coefficients {
		coefficients[i] = 0
	}
	return l
}

// UnpackedLength returns the number of elements that will be produced from
// unpacking a given binary.
func UnpackedLength(numElts, maxEltValue int) int {
	return Unpack(numElts, maxEltValue, nil, 0, nil, 0)
}

// Unpack unpacks a bit-packed array into an array of shorts.  The number of
// bits per element is implied by maxEltValue.
func Unpack(numElts, maxEltValue int, src []byte, srcOffset int, tgt []int16, tgtOffset int) int {
	bitsPerElement := int(CountBits(maxEltValue - 1))
	maxUsed := (numElts*bitsPerElement + 7) / 8
	if tgt == nil {
		return maxUsed
	}

	i, iMax := srcOffset, srcOffset+maxUsed-1
	j, jMax := tgtOffset, tgtOffset+numElts
	// tmp holds up to 16 bits from the source stream.
	// Stored as an int to make it easier to shift bits.
	tmp := int32(0xff & src[i])
	i++
	// tb holds the number of bits in tmp that are valid,
	// that is, that still need to be placed in the tgt array.
	// These will always be the least significant bits of tmp.
	tb, ob := 8, 0
	tgt[j] = 0

	_ = jMax // Original code doesn't use this.

	for i <= iMax || ob+tb >= bitsPerElement {
		if ob+tb < bitsPerElement {
			// Adding tb bits from tmp to the ob bits in tgt[j]
			// will not overflow the output element tgt[j].
			// Move all tb bits from tmp into tgt[j].
			shift := uint(bitsPerElement - ob - tb)
			tgt[j] |= int16((tmp << shift) & 0x0ffff)
			ob += tb
			tmp = int32((0xff & src[i]))
			i++
			tb = 8
		} else {
			// tmp has more bits than we need to finish output
			// element tgt[j]. Move some of the bits from tmp to
			// tgt[j] to finish it off, and save the leftovers in
			// tmp for the next iteration of the loop when we start
			// to fill in tgt[j+1].
			shift := uint(ob + tb - bitsPerElement)
			tgt[j] |= int16(((tmp & 0xff) >> shift) & 0xff)
			j++
			ob = 0
			tmp &= lowBitMask(shift)
			tb = int(shift)
		}
	}
	return maxUsed
}

// UnpackListedCoefficients unpacks a listed representation into a polynomial.
func UnpackListedCoefficients(f *poly.Polynomial, n, numOnes, numNegOnes int, in []byte, offset int) int {
	coefficients := make([]int16, numOnes+numNegOnes)
	l := Unpack(len(coefficients), n, in, offset, coefficients, 0)
	f.Reset()
	for i := 0; i < numOnes; i++ {
		f.Coeffs[coefficients[i]] = 1
	}
	for i := numOnes; i < len(coefficients); i++ {
		f.Coeffs[coefficients[i]] = -1
	}

	for i := range coefficients {
		coefficients[i] = 0
	}
	return l
}

// lowBitMask returns an integer that can be used to mask off the low numBits of
// a value.
func lowBitMask(numBits uint) int32 {
	return ^(-1 << uint(numBits))
}
