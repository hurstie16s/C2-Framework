package main

import (
	"Agent/internal/bootstrap"
	"Agent/pkg"
	"fmt"
)

func main() {
	if !(pkg.FileExists(pkg.GetCertPath()) && pkg.FileExists(pkg.GetKeyPath())) {
		fmt.Println("agent requires setup - bootstrapping")
		bootstrap.Init()
	}
	fmt.Println("agent ready - starting beacon")
}
