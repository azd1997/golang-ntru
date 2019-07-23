package polynomial

// Full is a polynomial represented by a slice of 16 bit signed coefficients.
type Polynomial struct {
	// P is the list of coefficients of the polynomial.
	P []int16
}

func (poly *Polynomial) IsInvertible() bool {

}