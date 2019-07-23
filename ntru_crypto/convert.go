package ntru_crypto


// convPolyBinaryToTrinaryHelper converts 3 bits to 2 trits.
func convPolyBinaryToTrinaryHelper(maxOffset, offset int, poly []int16, b int) {
	var a1, a2 int16
	switch b & 0x07 {
	case 0:
		a1, a2 = 0, 0
	case 1:
		a1, a2 = 0, 1
	case 2:
		a1, a2 = 0, -1
	case 3:
		a1, a2 = 1, 0
	case 4:
		a1, a2 = 1, 1
	case 5:
		a1, a2 = 1, -1
	case 6:
		a1, a2 = -1, 0
	case 7:
		a1, a2 = -1, 1
	}
	if offset < maxOffset {
		poly[offset] = a1
	}
	if offset+1 < maxOffset {
		poly[offset+1] = a2
	}
}

// convPolyBinaryToTrinaryHelper2 converts 24 bits stored in bits24 into 8
// trits.
func convPolyBinaryToTrinaryHelper2(maxOffset, offset int, poly []int16, bits24 int) {
	for i := 0; i < 24 && offset < maxOffset; i += 3 {
		shift := uint(24 - (i + 3))
		convPolyBinaryToTrinaryHelper(maxOffset, offset, poly, bits24>>shift)
		offset += 2
	}
}

// 将二项式转换为三项式。
// 输入为字节数组，看作是[0,1]的二项式，也即看做[]bit（bit-packed位组合）
// 输出为三项式系数数组切片 []int16，其中每一个元素只能为[-1,0,1]
func convPolyBinaryToTrinary(outputDegree int, bin []byte) []int16 {
	tri := make([]int16, outputDegree)
	blocks := len(bin) / 3
	remainder := len(bin) % 3

	// Perform the bulk of the conversion in 3-byte blocks.
	// 3 bytes == 24 bits --> 16 trits.
	for i := 0; i < blocks; i++ {
		val := int(bin[i*3])<<16 | int(bin[i*3+1])<<8 | int(bin[i*3+2])
		convPolyBinaryToTrinaryHelper2(outputDegree, 16*i, tri, val)
	}

	// Convert any partial block left at the end of the input buffer
	val := 0
	if remainder > 0 {
		val |= int(bin[blocks*3]) << 16
	}
	if remainder > 1 {
		val |= int(bin[blocks*3+1]) << 8
	}
	convPolyBinaryToTrinaryHelper2(outputDegree, 16*blocks, tri, val)

	return tri
}

// convPolyTritToBitHelper converts 2 trits to 3 bits, using the mapping
// defined in X9.92.
func convPolyTritToBitHelper(t1, t2 int16) byte {
	if t1 == -1 {
		t1 = 2
	}
	if t2 == -1 {
		t2 = 2
	}
	switch (t1 << 2) | t2 {
	case 0:
		return 0x00 // (t1,t2)=(  0,  0) ==> t = 0000
	case 1:
		return 0x01 // (t1,t2)=(  0,  1) ==> t = 0001
	case 2:
		return 0x02 // (t1,t2)=(  0, -1) ==> t = 0010
	case 4:
		return 0x03 // (t1,t2)=(  1,  0) ==> t = 0100
	case 5:
		return 0x04 // (t1,t2)=(  1,  1) ==> t = 0101
	case 6:
		return 0x05 // (t1,t2)=(  1, -1) ==> t = 0110
	case 8:
		return 0x06 // (t1,t2)=( -1,  0) ==> t = 1000
	case 9:
		return 0x07 // (t1,t2)=( -1,  1) ==> t = 1001
	default:
		return 0xff // (t1,t2)=( -1, -1) ==> t = 1010 (0xff)
	}
}

// convPolyTritToBitHelper2 converts 2 trits out of an array into a 3 bit value.
func convPolyTritToBitHelper2(offset int, trit []int16) byte {
	var t1, t2 int16
	if offset < len(trit) {
		t1 = trit[offset]
	}
	if offset+1 < len(trit) {
		t2 = trit[offset+1]
	}
	return convPolyTritToBitHelper(t1, t2)
}

// convPolyTrinaryToBinaryBlockHelper converts an array of 16 trits to 1 block
// (24 bits).
func convPolyTrinaryToBinaryBlockHelper(tOffset int, trit []int16, bOffset int, bits []byte) {
	a1 := int(convPolyTritToBitHelper2(tOffset, trit))
	a2 := int(convPolyTritToBitHelper2(tOffset+2, trit))
	a3 := int(convPolyTritToBitHelper2(tOffset+4, trit))
	a4 := int(convPolyTritToBitHelper2(tOffset+6, trit))
	a5 := int(convPolyTritToBitHelper2(tOffset+8, trit))
	a6 := int(convPolyTritToBitHelper2(tOffset+10, trit))
	a7 := int(convPolyTritToBitHelper2(tOffset+12, trit))
	a8 := int(convPolyTritToBitHelper2(tOffset+14, trit))

	// XXX: The ref calling code never checks this.
	// if (a1 | a2 | a3 | a4 | a5 | a6 | a7 | a8) == 0xff {
	//	return false
	// }

	// Pack the 8 3-bit values into a single 32-bit integer.
	// This makes it easier to pull off bytes later.
	val := a1<<21 | a2<<18 | a3<<15 | a4<<12 | a5<<9 | a6<<6 | a7<<3 | a8

	// Break the integer into bytes and put into the bits[] array.
	if bOffset < len(bits) {
		bits[bOffset] = byte(val >> 16)
		bOffset++
	}
	if bOffset < len(bits) {
		bits[bOffset] = byte(val >> 8)
		bOffset++
	}
	if bOffset < len(bits) {
		bits[bOffset] = byte(val)
	}
}

