//go:build !go1.16 && !go1.17
// +build !go1.16,!go1.17

package routine

import "errors"

var errUnsupported = errors.New("unsupported")

func getAllGoidByNative() (goids []int64, err error) {
	return nil, errUnsupported
}
