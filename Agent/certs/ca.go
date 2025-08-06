package certs

import (
	_ "embed"
)

//go:embed ca.crt
var caCertPEM []byte

func GetCACert() []byte {
	return caCertPEM
}
