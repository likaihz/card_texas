package xx

import (
	"log"
	"os"
)

var file *os.File

// create debug log
func Openlog(name string) {
	var err error
	file, err = os.Create(name + ".log")
	if err != nil {
		print("%s\r\n", err.Error())
		os.Exit(-1)
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("不知道这次能坚持几秒...")
}

// close debug log
func Closelog() {
	log.Println("擦，我又挂了！")
	err := file.Close()
	if err != nil {
		log.Println("Closelog() ", err)
	}
}
