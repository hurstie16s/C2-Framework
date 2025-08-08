package crypto

import (
	"Server/internal/handlers"
	"Server/pkg"
	"crypto/tls"
	"log"
	"net/http"
)

func TLSServer() {
	var tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{pkg.GetServerCert()},
		ClientCAs:    pkg.GetCAPool(),
		ClientAuth:   tls.RequireAnyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	var mux = http.NewServeMux()
	mux.HandleFunc("/bootstrap", handlers.BootstrapHandler)
	mux.HandleFunc("/beacon", requireMTLS(handlers.BeaconHandler))

	var server = &http.Server{
		Addr:      ":443",
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	log.Fatal(server.ListenAndServeTLS("", ""))

}

func requireMTLS(handler http.HandlerFunc) http.HandlerFunc {
	return nil
}
