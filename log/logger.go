package log

import (
	"log"
	"os"
)

var logger = log.New(
	os.Stderr,
	"service: ",
	//log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile,
	log.Ldate|log.Ltime|log.Lshortfile,
)
