package pkg

import (
	"net"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetMACAddress() []string {
	interfaces, _ := net.Interfaces()
	var macs []string
	for _, iface := range interfaces {
		mac := iface.HardwareAddr.String()
		if mac != "" {
			macs = append(macs, mac)
		}
	}
	return macs
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
