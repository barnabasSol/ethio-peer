package jwt

import (
	"os"
	"path/filepath"
)

func ReadPrivateKey() ([]byte, error) {
	absPath, err := filepath.Abs("internal/certs/private.pem")
	if err != nil {
		return nil, err
	}
	pk, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return pk, nil

}
