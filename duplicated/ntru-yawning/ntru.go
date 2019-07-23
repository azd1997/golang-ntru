package ntru_yawning

import (
	"github.com/azd1997/golang-ntru/duplicated"
	"github.com/azd1997/golang-ntru/duplicated/ntru-yawningtru-yawning/polynomial"
	"io"
)

// A PublicKey represents a NTRUEncrypt public key.
type PublicKey struct {
	Params *duplicated.KeyParams
	H      *duplicated.Polynomial
}


// A PrivateKey represents a NTRUEncrypt private key.
type PrivateKey struct {
	PublicKey
	F *duplicated.Polynomial
}

// GenerateKey generates a NTRUEncrypt keypair with the given parameter set
// using the random source random (for example, crypto/rand.Reader).
func GenerateKey(random io.Reader, oid duplicated.Oid) (priv *PrivateKey, err error) {
	keyParams := duplicated.Param(oid)
	if keyParams == nil {
		return nil, duplicated.InvalidParamError(oid)
	}
	prng := duplicated.readerToByteReader(random)
	igf := igf2.NewFromReader(keyParams.N, keyParams.C, prng)

	// Generate trinomial g that is invertible.
	var g *duplicated.Polynomial
	for isInvertible := false; !isInvertible; {
		if g, err = bpgm3.GenTrinomial(keyParams.N, keyParams.Dg+1, keyParams.Dg, igf); err != nil {
			return nil, err
		}
		gInv := inverterMod2048.Invert(g)
		isInvertible = gInv != nil
	}

	// Create F, f=1+p*F, and F^-1 mod q.
	var F, f, fInv *polynomial.Full
	for isInvertible := false; !isInvertible; {
		if F, err = bpgm3.GenTrinomial(keyParams.N, keyParams.Df, keyParams.Df, igf); err != nil {
			return nil, err
		}
		f = polynomial.New(int(keyParams.N))
		for i := range f.P {
			f.P[i] = (keyParams.P * F.P[i]) & 0xfff
		}
		f.P[0]++

		fInv = inverterMod2048.Invert(f)
		isInvertible = fInv != nil
	}

	// Calculate h = f^-1 * g * p mod q.
	h := polynomial.Convolution(fInv, g)
	for i := range h.P {
		h.P[i] = (h.P[i] * keyParams.P) % keyParams.Q
		if h.P[i] < 0 {
			h.P[i] += keyParams.Q
		}
	}

	fInv.Obliterate()
	F.Obliterate()

	priv = &PrivateKey{}
	priv.Params = keyParams
	priv.H = h
	priv.F = f
	return
}



