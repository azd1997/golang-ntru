package golang_ntru

//以多项式系数数组切片定义
type Polynomial []uint64;

//定义多项式环
type PolynomialRingN struct {
	degree uint64
	polynomial Polynomial
}

