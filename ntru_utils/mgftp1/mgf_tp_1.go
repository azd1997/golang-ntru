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

// 实现 MGF-TP-1 算法
// 将字节流转化为三进制位序列
// for converting a byte stream into a sequence of trits.
// It implements both the forward direction and the reverse.
package mgftp1

import (
	"io"

	"github.com/azd1997/golang-ntru/ntru_utils/poly"
)

// GenTrinomial generates a trinomial of degree N using the MGF-TP-1 algorithm
// to convert the io.ByteReader into trits.
func GenTrinomial(n int, r io.ByteReader) (*poly.Polynomial, error) {
	p := poly.NewPoly(n)

	limit := 5 * (n / 5)
	i := 0
	for i < limit {
		o, err := r.ReadByte()
		if err != nil {
			return nil, err
		} else if o >= 243 {
			continue
		}
		for j := 0; j < 5; j++ {
			b := o % 3
			p.Coeffs[i+j] = int16(b)
			o = (o - b) / 3
		}
		i += 5
	}
nLoop:
	for i < n {
		o, err := r.ReadByte()
		if err != nil {
			return nil, err
		} else if o >= 243 {
			continue
		}

		for j := 0; j < 5; j++ {
			b := o % 3
			p.Coeffs[i+j] = int16(b)
			o = (o - b) / 3
			if i+j+1 == n {
				break nLoop
			}
		}
		i += 5
	}

	// Renormalize from [0..2] to [-1..1]
	for i := range p.Coeffs {
		if p.Coeffs[i] == 2 {
			p.Coeffs[i] = -1
		}
	}
	return p, nil
}

// EncodeTrinomial generates a byte stream that is the encoding of a trinomial.
func EncodeTrinomial(poly *poly.Polynomial, w io.ByteWriter) error {
	n := len(poly.Coeffs)
	accum := byte(0)

	recenterTritTo0 := func(in int16) byte {
		if in == -1 {
			return 2
		}
		return byte(in)
	}

	// Encode 5 trits per byte, as long as we have >= 5 trits.
	for end := 5; end <= n; end += 5 {
		accum = recenterTritTo0(poly.Coeffs[end-1])
		accum = 3*accum + recenterTritTo0(poly.Coeffs[end-2])
		accum = 3*accum + recenterTritTo0(poly.Coeffs[end-3])
		accum = 3*accum + recenterTritTo0(poly.Coeffs[end-4])
		accum = 3*accum + recenterTritTo0(poly.Coeffs[end-5])
		if err := w.WriteByte(accum); err != nil {
			return err
		}
	}
	if end := n - (n % 5); end < n {
		n--
		accum = recenterTritTo0(poly.Coeffs[n])
		for end < n {
			n--
			accum = 3*accum + recenterTritTo0(poly.Coeffs[n])
		}
		if err := w.WriteByte(accum); err != nil {
			return err
		}
	}
	return nil
}
