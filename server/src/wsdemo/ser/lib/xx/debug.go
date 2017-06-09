package xx

import (
	"fmt"
	"log"
	"os"
)

func Log2file(path string) (*log.Logger, *os.File) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	l := log.New(file, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	return l, file
}
