package main

import (
	"bytes"
	"log"
	"os"
)

func onError(err error, messages ...string) {
	if err != nil {

		var buffer bytes.Buffer

		for index, message := range messages {
			if index == 0 {
				buffer.WriteString("[ERROR] ")
			}
			buffer.WriteString(message)
			buffer.WriteString(" ")
		}

		log.Fatal(buffer.String())
		os.Exit(1)
	}
}
