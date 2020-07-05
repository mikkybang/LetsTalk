package model

import (
	"log"
	"os"
)

func dropboxFileUploader(file *os.File) {
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("unable to close file", err)
		}

		if err := os.Remove(file.Name()); err != nil {
			log.Println("unable to remove file", err)
		}
	}()
}
