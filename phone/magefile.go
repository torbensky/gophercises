// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	return sh.Run("go", "build", "-o", "phone")
}

func Postgres() error {
	fmt.Println("Staring postgres in docker...")
	return sh.Run("docker", "run", "--rm", "-p", "5432:5432", "--name", "test-postgres", "-e", "POSTGRES_PASSWORD=gophone", "-d", "postgres")
}

func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return sh.Run("dep", "ensure")
}

func Clean() error {
	fmt.Println("Cleaning...")
	os.RemoveAll("phone")
	return sh.Run("docker", "stop", "test-postgres")
}
