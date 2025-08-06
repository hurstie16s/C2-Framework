package pkg

const (
	certPath = "certs/client.crt"
	keyPath  = "certs/client.key"
)

// Getters

func GetCertPath() string {
	return certPath
}

func GetKeyPath() string {
	return keyPath
}
