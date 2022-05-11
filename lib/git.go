package lib

import (
	"fmt"
	"strings"

	"github.com/pwuersch/ggt/lib/globals"
)

func ParseGitDestDir(url string) string {
	parts := strings.Split(url, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}

func Clone(url string, dest string, profile globals.Profile) error {
	err := NewRunner("git", "clone", url, dest).Run()
	if err != nil {
		return err
	}

	return Config(dest, profile)
}

func Config(dir string, profile globals.Profile) error {
	commitName := profile.CommitName

	fmt.Printf("Setting git user.name to %s\n", commitName)
	err := NewRunner("git", "config", "user.name", commitName).WithDir(dir).Run()
	if err != nil {
		return err
	}

	fmt.Printf("Setting git user.email to %s\n", profile.Email)
	return NewRunner("git", "config", "user.email", profile.Email).WithDir(dir).Run()
}

func Add() error {
	return NewRunner("git", "add", ".").Run()
}

func Commit(message string) error {
	return NewRunner("git", "commit", "-m", message).Run()
}

func Push() error {
	return NewRunner("git", "push").Run()
}
