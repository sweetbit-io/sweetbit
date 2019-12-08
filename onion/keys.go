package onion

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/go-errors/errors"
)

type Version int

const (
	V2 Version = iota
	V3  // TODO(davidknezic): support V3
)

func GeneratePrivateKey(v Version) (*rsa.PrivateKey, error) {
	if v != V2 {
		return nil, errors.New("only V2 supported for now")
	}

	//switch v {
	//case V2:
	//	return control.GenKey(control.KeyAlgoRSA1024)
	//case V3:
	//	return control.GenKey(control.KeyAlgoED25519V3)
	//}

	// Generate a V2 RSA 1024 bit key
	return rsa.GenerateKey(rand.Reader, 1024)

	//return ed25519.GenerateKey(rand.Reader)
}
