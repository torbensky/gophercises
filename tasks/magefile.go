// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
)

func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "tasks", ".")
	return cmd.Run()
}

func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("dep", "ensure")
	return cmd.Run()
}

func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("tasks")
}
