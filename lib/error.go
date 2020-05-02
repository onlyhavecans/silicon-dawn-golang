package lib

import "log"

func FatalIfErr(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
