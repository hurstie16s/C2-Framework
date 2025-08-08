package pkg

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

var (
	serverCert tls.Certificate
	caCert     []byte
	caKey      []byte
	caPool     *x509.CertPool
	err        error
)

// exposed functions

func LoadServerCert(pathToCert string, pathToKey string) {
	serverCert, err = tls.LoadX509KeyPair(pathToCert, pathToKey)
	CheckError(err)
}

func LoadCACert(pathToCert string, pathToKey string) {
	caCert, err = os.ReadFile(pathToCert)
	CheckError(err)

	caKey, err = os.ReadFile(pathToKey)
	CheckError(err)

	caPool = x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)
}

func GetCACert() []byte {
	return caCert
}

func GetCAKey() []byte {
	return caKey
}

func GetServerCert() tls.Certificate {
	return serverCert
}

func GetCAPool() *x509.CertPool {
	return caPool
}
