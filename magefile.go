//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	binName = "git-trash"
	mainPkg = "./cmd"
	outDir  = "./dist"

	moduleName = "github.com/ajm113/git-trash"
)

// Default runs Build
var Default = Build

func Build() error {
	mg.Deps(Tidy)
	fmt.Printf("Building %s...\n", binName)
	return sh.RunV("go", "build",
		"-ldflags", ldflags(),
		"-o", filepath.Join(outDir, binName),
		mainPkg,
	)
}

func Run() error {
	mg.Deps(Build)
	bin := filepath.Join(outDir, binName)
	args := os.Args[1:] // pass through any extra args
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Tidy runs go mod tidy
func Tidy() error {
	fmt.Println("Tidying modules...")
	return sh.RunV("go", "mod", "tidy")
}

func Clean() error {
	fmt.Println("Cleaning...")
	paths := []string{outDir}
	for _, p := range paths {
		if err := os.RemoveAll(p); err != nil {
			return err
		}
	}
	return nil
}

func Test() error {
	fmt.Println("Running tests...")
	return sh.RunV("go", "test", "-v", "-race", "-count=1", "./...")
}

func ldflags() string {
	version := gitVersion()
	commit := gitCommit()
	date := time.Now().UTC().Format(time.RFC3339)
	base := fmt.Sprintf("%s/internal/version", moduleName)

	return fmt.Sprintf(
		`-s -w -X %s.Version=%s -X %s.Commit=%s -X %s.Date=%s`,
		base, version,
		base, commit,
		base, date,
	)
}

func gitVersion() string {
	v, _ := sh.Output("git", "describe", "--tags", "--always", "--dirty")
	if v == "" {
		return "dev"
	}
	return v
}

func gitCommit() string {
	c, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	if c == "" {
		return "unknown"
	}
	return c
}
