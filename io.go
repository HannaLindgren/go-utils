package goutils

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ReadFileToLines Read a file into a list of lines
func ReadFileToLines(fName string) ([]string, error) {
	s, err := ReadFileToString(fName)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(strings.TrimSuffix(s, "\n"), "\n"), nil
}

// ReadFileToString Read a file into a string using ioutil.ReadFile
func ReadFileToString(fName string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Clean(fName))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadStdinToString Read stdin into a string using ioutil.ReadFile
func ReadStdinToString() (string, error) {
	stdin := bufio.NewReader(os.Stdin)
	b, err := ioutil.ReadAll(stdin)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
