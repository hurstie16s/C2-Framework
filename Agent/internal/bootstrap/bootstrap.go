package bootstrap

import (
	"Agent/certs"
	"Agent/pkg"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func Init() {
	fmt.Println("bootstrapping agent")

	// generate private key
	privKey, _ := rsa.GenerateKey(rand.Reader, 4096)

	// get device MAC addresses and take first
	var mac = pkg.GetMACAddress()[0]
	var agentID = "agent-" + mac

	// create csr
	csrPEM := createCSR(agentID, privKey)

	// setup HTTPS client with server
	client := setupBootstrapComms()

	// send csr to server
	// TODO: check what url should be
	resp := pkg.Post(
		client,
		"https://c2server.com/bootstrap",
		"application/pem-certificate-chain",
		csrPEM,
	)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			pkg.CheckError(err)
		}
	}(resp.Body)

	// save crt and key
	saveCertKey(resp, privKey)

	fmt.Println("agent bootstrapped")
	return
}

func createCSR(agentID string, privKey *rsa.PrivateKey) []byte {
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: agentID,
		},
	}
	csrDER, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privKey)
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})

	return csrPEM
}

func setupBootstrapComms() *http.Client {
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(certs.GetCACert()) {
		// TODO: find a way to recover
		panic("failed to append root certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs: caPool,
	}
	transport := http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: &transport, Timeout: time.Second * 10}

	return client
}

func saveCertKey(resp *http.Response, privKey *rsa.PrivateKey) {
	certPEM, err := io.ReadAll(resp.Body)
	pkg.CheckError(err)
	// TODO: ensure going to correct location
	err = os.WriteFile(pkg.GetCertPath(), certPEM, 0644)
	pkg.CheckError(err)

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	)
	err = os.WriteFile(pkg.GetKeyPath(), keyPEM, 0600)
	pkg.CheckError(err)
}
