package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

const (
	cardsDir = "data"
	cardsURL = "http://egypt.urnash.com/media/blogs.dir/1/files/2018/01/The-Tarot-of-the-Silicon-Dawn.zip"
	exitFail = 1
)

func main() {
	if err := run(os.Args, cardsDir, cardsURL, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(exitFail)
	}
}

func run(_ []string, outDir string, url string, stdout io.Writer) error {
	fmt.Fprintf(stdout, "Creating card directory %s\n", outDir)
	err := os.MkdirAll(outDir, 0o700)
	if err != nil {
		fmt.Fprintf(stdout, "os.MkDirAll err = %v\n", err)
	}

	fmt.Fprint(stdout, "Getting Zip File\n")
	z, err := retrieveZip(url)
	if err != nil {
		err = fmt.Errorf("retrieveZip(%s) err = %w", url, err)
		return err
	}

	fmt.Print("Unzipping files\n")
	err = unzipFiles(z, outDir)
	if err != nil {
		err = fmt.Errorf("unzipZiles(zipfile, %q) err = %w", outDir, err)
		return err
	}

	fmt.Fprint(stdout, "Finished!\n")
	return nil
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

		if skipFile(f.Name()) {
			_ = f.Close()
			continue
		}

		fileName := filepath.Join(destinationDir, f.Name())
		body, err := io.ReadAll(f)
		if err != nil {
			fmt.Print("Failure to read compressed file ", f.Name(), err)
			_ = f.Close()
			continue
		}

		err = os.WriteFile(fileName, body, 0o644)
		if err != nil {
			_ = f.Close()
			return err
		}

		_ = f.Close()
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
	return io.ReadAll(resp.Body)
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
