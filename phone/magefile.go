// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const pgUser = "postgres"
const pgPass = "gophone"
const pgDbName = "postgres"

func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	return sh.Run("go", "build", "-o", "phone")
}

func Run() error {
	mg.Deps(Build)
	fmt.Println("Running...")
	return sh.RunWith(
		map[string]string{
			"PG_DB_NAME":  pgDbName,
			"PG_PASSWORD": pgPass,
			"PG_USER":     pgUser,
		},
		"./phone",
	)
}

func SeedDatabase() error {
	return sh.RunWith(
		map[string]string{
			"PGPASSWORD": pgPass,
		},
		"psql", "-h", "localhost", "-U", pgUser, "-d", pgDbName, "-f", "fixture.sql",
	)
}

func Postgres() error {
	fmt.Println("Starting postgres in docker...")
	return sh.Run("docker", "run", "--rm", "-p", "5432:5432", "--name", "test-postgres", "-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", pgPass), "-d", "postgres")
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
