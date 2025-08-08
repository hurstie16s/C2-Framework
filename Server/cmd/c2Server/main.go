package main

import (
	"Server/internal/crypto"
	"Server/pkg"
)

func main() {

	// load certificates
	pkg.LoadServerCert(
		"certs/server.crt",
		"certs/server.key",
	)
	pkg.LoadCACert(
		"certs/ca.crt",
		"certs/ca.key",
	)

	// start TLS server
	crypto.TLSServer()
}
