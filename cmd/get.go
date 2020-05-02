/*
Copyright © 2020 Amelia Aronsohn <squirrel@wearing.black>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"onlyhavecans.works/amy/silicondawn/lib"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download and prepare the silicon-dawn",
	Long: `We all have our own methods for preparing our space and practice.
Before we can draw card for any one we must first acquire it.

There is no local show which sells it, so instead we take its digital format.
Luckily for us it is offered for free by it's creator. This command is for
acquiring this package, opening it, and laying out our cards`,
	Run: func(cmd *cobra.Command, args []string) {
		cardsURL := viper.GetString("CardsURL")
		cardsDirectory := viper.GetString("CardsDirectory")

		err := os.MkdirAll(cardsDirectory, 0700)
		lib.FatalIfErr("Making directory", err)

		log.Print("Getting Zip File")
		z, err := retrieveZip(cardsURL)
		lib.FatalIfErr("Download File Failed", err)

		log.Print("Unzipping files")
		err = unzipFiles(z, cardsDirectory)
		lib.FatalIfErr("Unzip failed", err)

		log.Print("Finished!")
	},
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
			f.Close()
			continue
		}

		fileName := filepath.Join(destinationDir, f.Name())
		body, err := ioutil.ReadAll(f)
		if err != nil {
			log.Print("Failure to read compressed file ", f.Name(), err)
			f.Close()
			continue
		}

		err = ioutil.WriteFile(fileName, body, 0644)
		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			log.Print("close failed", err)
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
	default:
		return false
	}
}

func init() {
	rootCmd.AddCommand(getCmd)
}
