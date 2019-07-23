package ntru_yawning

import (
	"fmt"
	"github.com/azd1997/golang-ntru/duplicated"
)

// InvalidParamError is the error returned when an invalid parameter set is
// specified during key generation.
type InvalidParamError duplicated.Oid

func (e InvalidParamError) Error() string {
	return fmt.Sprintf("ntru: unsupported OID: %d", e)
}
