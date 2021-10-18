package integration_tests

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func writeFile(path string, body []byte) error {
	_, err := os.Create(path)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, body, 0600)
}
