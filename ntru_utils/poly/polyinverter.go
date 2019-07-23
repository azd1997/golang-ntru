package poly

type Inverter interface {
	//多项式求逆接口
	Invert(p *Polynomial) *Polynomial
}

// 实现在环(Z/pZ)[X]/(X^N-1)上为多项式寻找模p逆.(这里p表示质数)
type InverterModPrime struct {
	// 模数(modulus)要求为质数
	prime int

	// invModPrime is a table of inverses mod prime, setup so that
	// invModPrime[i] * i = 1 (mod prime) if the inverse of i exists, and
	// invModPrime[i] = 0 if thte inverse of i does not exist.
	invModPrime []int16
}

/*对系数x进行模p处理，使系数限制在[0,1,2,...,prime-1]*/
func (v *InverterModPrime) modPrime(x int16) int16 {
	ret := int(x) % v.prime
	if ret < 0 {
		ret += v.prime
	}
	return int16(ret)
}

/*在环(Z/pZ)[X]/(X^N-1)上求多项式的逆.*/
// 参阅 NTRU Cryptosystems Tech Report #014 "Almost Inverses and Fast NTRU Key
// Creation."
func (v *InverterModPrime) Invert(a *Polynomial) *Polynomial {
	//待求逆多项式可能的最高阶次N-1
	N := len(a.Coeffs)

	// 初始化: (bcfg均为最高N阶的多项式)
	// k = 0, B(X) = 1, C(X) = 0, f(X)=a(X), g(X)=X^N-1
	k := 0
	b := NewPoly(N + 1)
	c := NewPoly(N + 1)
	f := NewPoly(N + 1)
	g := NewPoly(N + 1)
	b.Coeffs[0] = 1
	for i := 0; i < N; i++ {
		f.Coeffs[i] = v.modPrime(a.Coeffs[i])
		//f.Coeffs[i] = a.Coeffs[i]
	}
	g.Coeffs[N] = 1
	g.Coeffs[0] = int16(v.prime - 1)  //int16(v.prime - 1)	//TODO: 这里 prime=0 ?

	// 记录 f,g 阶次
	df := f.getDegree()
	dg := N

	//主循环
	for {

		// while f[0] = 0 {f/=X, c*=X, k++}
		// 直到f(x)常数项不为0，x项可以不存在
		// 在这里c(x)*x并没有变化，当c不为0以后，再返回此处会有变化
		for f.Coeffs[0] == 0 && df > 0 {
			df--
			f.divideByX()
			c.multiplyByX()
			k++
		}

		if df == 0 {
			// 当f只含常数项时要确保该常数项有逆存在。
			// 例子：p=3时，f[0]= +-1而不能能为0，就是因为为0时没有逆
			// Make sure there is a solution, return nil if a is not invertible.
			f0Inv := v.invModPrime[f.Coeffs[0]]
			if f0Inv == 0 {
				return nil
			}

			// b(X) = f[0]inv * b(X) mod p
			// return X^(N-k) * b
			//shift := (N - k) % N	//注意：根据主循环第一次循环中的第一步，k确实一定小于N,但是继续循环下去，k是会超过N的
			shift := N-k
			shift %= N

			if shift < N {
				shift += N	//确保shift在[0,2N]
			}
			ret := NewPoly(N)
			for i := range ret.Coeffs {
				// b(X) = X^(N-k)*f[0]inv * b(X) mod p
				ret.Coeffs[(i+shift)%N] = v.modPrime(f0Inv * b.Coeffs[i])
			}
			return ret
		}

		if df < dg {
			// swap(f,g), swap(b,c), swap(df, dg)
			f, g = g, f
			b, c = c, b
			df, dg = dg, df
		}

		// u = f[0] * g[0]inv mod p
		u := v.modPrime(f.Coeffs[0] * v.invModPrime[g.Coeffs[0]])

		// f(X) -= u*g(X) mod p
		for i := range f.Coeffs {
			f.Coeffs[i] = v.modPrime(f.Coeffs[i] - u*g.Coeffs[i])
		}

		// b(X) -= u*c(X) mod p
		for i := range b.Coeffs {
			b.Coeffs[i] = v.modPrime(b.Coeffs[i] - u*c.Coeffs[i])
		}
	}
}

/*新建模P逆对象*/
func NewInverterModPrime(prime int, invModPrime []int16) *InverterModPrime {
	return &InverterModPrime{prime: prime, invModPrime: invModPrime}
}

/*求模素数的指数幂的逆*/
//用途，求模q逆。q一般为2的指数幂
type InverterModPowerOfPrime struct {
	primeInv *InverterModPrime

	powerOfPrime int16
}

/*在环(Z/p^rZ)[X]/(X^N-1)上求多项式的逆*/
func (v *InverterModPowerOfPrime) Invert(a *Polynomial) *Polynomial {
	// b = a inverse mod prime.
	b := v.primeInv.Invert(a)
	if b == nil {
		return nil
	}

	//for q := int(v.primeInv.prime); q < int(v.powerOfPrime); q *= q {}
	for q := int(v.primeInv.prime); q < int(v.powerOfPrime); {
		q *= q	//注意！q*=q不要放到循环的第三句去!原因是放到循环第三句是第一次循环结束才会执行q*=q
		// b(X) = b(X) * (2-a(X)b(X)) (mod q)
		//  i : c = a*b
		c := ConvolutionModN(a, b, q)
		// ii : c = 2-a*b
		c.Coeffs[0] = 2 - c.Coeffs[0]
		if c.Coeffs[0] < 0 {
			c.Coeffs[0] += int16(q)
		}
		for i := 1; i < len(b.Coeffs); i++ {
			c.Coeffs[i] = int16(q - int(c.Coeffs[i])) // This is -c (mod q)
		}
		b = ConvolutionModN(b, c, q)
	}
	return b
}


func NewInverterModPowerOfPrime(powerOfPrime int16, prime int, invModPrime []int16) *InverterModPowerOfPrime {
	v := &InverterModPowerOfPrime{powerOfPrime: powerOfPrime}
	v.primeInv = NewInverterModPrime(prime, invModPrime)
	return v
}


var _ Inverter = (*InverterModPrime)(nil)
var _ Inverter = (*InverterModPowerOfPrime)(nil)