package ntru_crypto

import (
	"errors"
	"fmt"
	"github.com/azd1997/golang-ntru/ntru_utils/params"
)

// InvalidParamError is the error returned when an invalid parameter set is
// specified during key generation.
type InvalidParamError params.Oid

func (e InvalidParamError) Error() string {
	return fmt.Sprintf("ntru: unsupported OID: %d", e)
}

// ErrMessageTooLong is the error returned when a message is too long for the
// parameter set.
var ErrMessageTooLong = errors.New("ntru: message too long for chosen parameter set")

// ErrDecryption is the error returned when decryption fails.  It is
// deliberately vague to avoid adaptive attacks.
var ErrDecryption = errors.New("ntru: decryption error")

// ErrInvalidKey is the error returned when the key is invalid.
var ErrInvalidKey = errors.New("ntru: invalid key")

