package poly

type Polynomial struct {
	Coeffs []int16
}

/*获取多项式的阶次*/
func (poly *Polynomial) getDegree() (deg int) {
	deg = len(poly.Coeffs) - 1
	for deg>0 && poly.Coeffs[deg]==0 {
		deg--
	}
	return
}

/*多项式除以x，即f(x)/x*/
func (poly *Polynomial) divideByX() {
	f0 := poly.Coeffs[0]
	for i := 0; i < len(poly.Coeffs)-1; i++ {
		poly.Coeffs[i] = poly.Coeffs[i+1]
	}
	poly.Coeffs[len(poly.Coeffs)-1] = f0
}

/*多项式乘以x，即f(x)*x*/
func (poly *Polynomial) multiplyByX() {
	fLast := poly.Coeffs[len(poly.Coeffs)-1]
	for i:=len(poly.Coeffs)-1;i>0;i-- {
		poly.Coeffs[i] = poly.Coeffs[i-1]
	}
	poly.Coeffs[0] = fLast
}

/*对多项式系数进行模q处理，使得其系数重定位到[newLowerLimit, newLowerLimit+q]
这使得旧的系数作模q处理等于新的系数*/
//若选择模数为3，下限为-1，则可将多项式系数约束到[-1,0,1]
func (poly *Polynomial) recenterModQ(q,newLowerLimit int) {
	newUpperLimit := newLowerLimit + q
	for i := range poly.Coeffs {
		tmp := int(poly.Coeffs[i]) % q
		if tmp >= newUpperLimit {
			tmp -= q
		}
		if tmp < newLowerLimit {
			tmp += q
		}
		poly.Coeffs[i] = int16(tmp)
	}
}

/***********************************多项式的新建************************************/

/*新建n-1阶多项式*/
func NewPoly(n int) (p *Polynomial) {
	p = &Polynomial{}
	p.Coeffs = make([]int16, n)
	return
}

/*根据系数数组创建多项式*/
func NewFromCoeffs(coeffs []int16) (p *Polynomial) {
	p = NewPoly(len(coeffs))
	copy(p.Coeffs, coeffs)
	return
}


/***********************************多项式的乘法************************************/

/*对两个N阶多项式进行环上卷积运算，也叫星乘计算*/
func Convolution(a, b *Polynomial) (c *Polynomial) {
	if len(a.Coeffs) != len(b.Coeffs) {
		// XXX: Does this happen ever?
		c = NewPoly(0)
	} else {
		c = NewPoly(len(a.Coeffs))
		for i := range a.Coeffs {
			for j := range b.Coeffs {
				c.Coeffs[(i+j)%len(c.Coeffs)] += (a.Coeffs[i] * b.Coeffs[j])
			}
		}
	}
	return
}

/*计算两多项式卷积，并将系数限制在[0,..,Q]*/
//例如coefficientModulus=3 => 系数限制在[0,1,2]
func ConvolutionModN(a, b *Polynomial, coefficientModulus int) (c *Polynomial) {
	c = Convolution(a, b)
	c.recenterModQ(coefficientModulus, 0)
	return
}


/***********************************多项式的加法************************************/

/*与另一多项式相加，结果多项式的系数再进行重定位*/
func (a *Polynomial) AddAndRecenter(b *Polynomial, coefficientModulus, newLowerLimit int) (c *Polynomial) {
	c = NewPoly(len(a.Coeffs))
	for i := range c.Coeffs {
		c.Coeffs[i] = a.Coeffs[i] + b.Coeffs[i]
	}
	c.recenterModQ(coefficientModulus, newLowerLimit)
	return
}

/*与另一多项式相加，结果多项式的系数再进行重定位,限制在[0,..,coefficientModulus-1]*/
func (a *Polynomial) Add(b *Polynomial, coefficientModulus int) (c *Polynomial) {
	return a.AddAndRecenter(b, coefficientModulus, 0)
}

/***********************************多项式的减法************************************/

/*与另一多项式相加，结果多项式的系数再进行重定位*/
func (a *Polynomial) SubtractAndRecenter(b *Polynomial, coefficientModulus, newLowerLimit int) (c *Polynomial) {
	c = NewPoly(len(a.Coeffs))
	for i := range c.Coeffs {
		c.Coeffs[i] = a.Coeffs[i] - b.Coeffs[i]
	}
	c.recenterModQ(coefficientModulus, newLowerLimit)
	return
}

/*与另一多项式相加，结果多项式的系数再进行重定位,限制在[0,..,coefficientModulus-1]*/
func (a *Polynomial) Subtract(b *Polynomial, coefficientModulus int) (c *Polynomial) {
	return a.SubtractAndRecenter(b, coefficientModulus, 0)
}

/***********************************多项式的比较************************************/

/*与另一多项式进行比较，看是否相等*/
func (a *Polynomial) Equals(b *Polynomial) bool {
	//先比较长度是否相等
	if len(a.Coeffs) != len(b.Coeffs) {
		return false
	}
	//比较每一项系数是否相等
	res := int16(0)
	for i := range a.Coeffs {
		res |= a.Coeffs[i] ^ b.Coeffs[i]	//"|"按位或，"^"按位异或。如果a,b对应系数的二进制位有一个不相等，res就不可能为0
	}
	return res == 0
}

/***********************************多项式的清零重置************************************/

/*多项式系数全置为0*/
func (a *Polynomial) Reset() {
	for i := range a.Coeffs {
		a.Coeffs[i] = 0
	}
}



