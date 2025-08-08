package handlers

import (
	"Server/pkg"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"log"
	"net/http"
	"time"
)

func BootstrapHandler(w http.ResponseWriter, r *http.Request) {

	// check request
	var req *pkg.BootstrapRequest
	var csr *x509.CertificateRequest
	var ok bool
	req, ok = checkRequest(http.MethodPost, w, r)
	if !ok {
		return
	}

	// check CSR
	csr, ok = checkCSR(req, w)
	if !ok {
		return
	}

	// signCSR
	signCSR(w, csr)
}

func checkRequest(
	method string,
	w http.ResponseWriter,
	r *http.Request,
) (*pkg.BootstrapRequest, bool) {

	var req *pkg.BootstrapRequest

	if r.Method != method {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return nil, false
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, false
	}

	return req, true
}

func checkCSR(
	req *pkg.BootstrapRequest,
	w http.ResponseWriter,
) (*x509.CertificateRequest, bool) {

	block, _ := pem.Decode([]byte(req.CSR))
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		http.Error(w, "invalid CSR", http.StatusBadRequest)
		return nil, false
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		http.Error(w, "failed to parse CSR", http.StatusBadRequest)
		return nil, false
	}

	if err := csr.CheckSignature(); err != nil {
		http.Error(w, "invalid CSR Signature", http.StatusBadRequest)
		return nil, false
	}

	return csr, true
}

func signCSR(
	w http.ResponseWriter,
	csr *x509.CertificateRequest,
) {

	caCertBlock, _ := pem.Decode(pkg.GetCACert())
	caKeyBlock, _ := pem.Decode(pkg.GetCAKey())

	caCert, _ := x509.ParseCertificate(caCertBlock.Bytes)
	caKey, _ := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)

	clientCertDER, err := x509.CreateCertificate(
		nil,
		&x509.Certificate{
			SerialNumber: pkg.BigSerial(),
			Subject:      csr.Subject,
			NotBefore:    time.Now(),
			NotAfter:     time.Now().Add(365 * 24 * time.Hour),
			KeyUsage: x509.KeyUsageDigitalSignature |
				x509.KeyUsageKeyEncipherment,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			BasicConstraintsValid: true,
		},
		caCert,
		csr.PublicKey,
		caKey,
	)
	if err != nil {
		http.Error(
			w,
			"could not sign certificate",
			http.StatusInternalServerError,
		)
		return
	}

	pem.Encode(
		w,
		&pem.Block{Type: "CERTIFICATE", Bytes: clientCertDER},
	)

	log.Println("Agent Bootstrapped")
}
