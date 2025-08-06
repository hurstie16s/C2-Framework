package main

import (
	"Agent/pkg"
	"fmt"
)

func main() {
	if pkg.FileExists(pkg.GetCertPath()) && pkg.FileExists(pkg.GetKeyPath()) {
		fmt.Println("agent setup - skipping bootstrap")
	} else {
		fmt.Println("agent requires setup - bootstrapping")
	}
}
