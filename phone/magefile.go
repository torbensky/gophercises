// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/sh"
)

func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	return sh.Run("go", "build")
}

func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return sh.Run("dep", "ensure")
}

func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}
