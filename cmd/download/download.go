package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

func main() {
	cardsDirectory := "data"
	cardsURL := "http://egypt.urnash.com/media/blogs.dir/1/files/2018/01/The-Tarot-of-the-Silicon-Dawn.zip"

	fmt.Printf("Creating card directory %s\n", cardsDirectory)
	err := os.MkdirAll(cardsDirectory, 0700)
	if err != nil {
		fmt.Printf("Making Directory: %w\n", err)
	}

	fmt.Print("Getting Zip File\n")
	z, err := retrieveZip(cardsURL)
	if err != nil {
		fmt.Printf("Downloading file: %w\n", err)
	}

	fmt.Print("Unzipping files\n")
	err = unzipFiles(z, cardsDirectory)
	if err != nil {
		fmt.Printf("Unzipping file: %w\n", err)
	}

	fmt.Print("Finished!\n")
}

func unzipFiles(zipData []byte, destinationDir string) error {
	// Zip Files need to know the size of the file
	contentLen := len(zipData)
	r := bytes.NewReader(zipData)
	z := archiver.NewZip()

	err := z.Open(r, int64(contentLen))
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer z.Close()

	// iterate each file in the archive until EOF
	for {
		f, err := z.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		//noinspection GoUnhandledErrorResult
		defer f.Close()

		if skipFile(f.Name()) {
			continue
		}

		fileName := filepath.Join(destinationDir, f.Name())
		body, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Print("Failure to read compressed file ", f.Name(), err)
			continue
		}

		err = ioutil.WriteFile(fileName, body, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func retrieveZip(url string) ([]byte, error) {
	// Basic HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func skipFile(name string) bool {
	switch {
	case strings.HasPrefix(name, "._"):
		return true
	case strings.HasPrefix(name, "__MACOSX"):
		return true
	case strings.HasPrefix(name, "sand-home"):
		return true
	case strings.HasSuffix(name, "jpg"):
		return false
	case strings.HasSuffix(name, "png"):
		return false
	default:
		return true
	}
}
