package io

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsDirectory(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return false, err
	}
	file.Close()
	return fileInfo.IsDir(), nil
}

func DetachFileExtension(fName string) (string, string) {
	ext := filepath.Ext(fName)
	name := fName[:len(fName)-len(ext)]
	return name, strings.TrimPrefix(ext, ".")
}

func RemoveFileExtension(fName string) string {
	res, _ := DetachFileExtension(fName)
	return res
}

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

// ReadStdinToLines Read stdin into a list of lines
func ReadStdinToLines() ([]string, error) {
	s, err := ReadStdinToString()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(strings.TrimSuffix(s, "\n"), "\n"), nil
}

// IsFile returns true if the given file exists (as a file or as a directory)
func IsFile(fName string) bool {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetFileReader reads an input file, gzipped or plain text, and returns an io.Reader for line scanning, along with the file handle, that needs to be closed after reading.
func GetFileReader(fName string) (io.Reader, *os.File, error) {
	fh, err := os.Open(filepath.Clean(fName))
	//defer fh.Close()
	if err != nil {
		return nil, fh, fmt.Errorf("couldn't open file %s for reading : %v", fName, err)
	}

	if strings.HasSuffix(fName, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			return nil, fh, fmt.Errorf("couldn't to open gz reader : %v", err)
		}
		return io.Reader(gz), fh, nil
	}
	return io.Reader(fh), fh, nil
}
