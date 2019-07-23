package poly

// 计算多项式的模4结果并返回为 bit-packed byte array.
func CalcPolyMod4Packed(r *Polynomial) (r4 []byte) {
	// R4 ： 2 bits per element, 4 elements per byte.
	r4 = make([]byte, (len(r.Coeffs)+3)/4)	//+3确保/4后有足够多byte来表示

	var i, j int
	for ; i < len(r4)-1; i, j = i+1, j+4 {
		tmp := (r.Coeffs[j] & 0x03) << 6
		tmp |= (r.Coeffs[j+1] & 0x03) << 4
		tmp |= (r.Coeffs[j+2] & 0x03) << 2
		tmp |= r.Coeffs[j+3] & 0x03
		r4[i] = byte(tmp)
	}

	remElements := len(r.Coeffs) & 3
	if remElements > 0 {
		r4[i] |= byte(r.Coeffs[j]&0x03) << 6
	}
	if remElements > 1 {
		r4[i] |= byte(r.Coeffs[j+1]&0x03) << 4
	}
	if remElements > 2 {
		r4[i] |= byte(r.Coeffs[j+2]&0x03) << 2
	}
	return
}

// 检查多项式中系数-1,0,1个数是否大于一定值
func CheckDm0(p *Polynomial, dm0 int16) bool {
	var numOnes, numNegOnes int16
	for _, v := range p.Coeffs {
		switch v {
		case -1:
			numNegOnes++
		case 1:
			numOnes++
		}
	}
	if numOnes < dm0 || numNegOnes < dm0 || int16(len(p.Coeffs))-(numOnes+numNegOnes) < dm0 {
		return false
	}
	return true
}
