package logfile

import (
	"log"
	"os"
)

func Init(filename string) (err error, logPanicAndCloseFile func()) {
	// create if doesnt exist, truncate (delete all contents) if already exists
	file, err := os.Create(filename)
	if err != nil {
		return err, nil
	}

	log.SetOutput(file)
	logPanicAndCloseFile = func() {
		defer file.Close()
		if r := recover(); r != nil {
			log.Println("logging panic before file close")
			log.Panicf("%+v", r)
		}
	}
	return nil, logPanicAndCloseFile
}
