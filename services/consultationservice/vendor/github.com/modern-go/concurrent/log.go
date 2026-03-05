package concurrent

import (
	"os"
	"log"
	"io/ioutil"
)


var ErrorLogger = log.New(os.Stderr, "", 0)


var InfoLogger = log.New(ioutil.Discard, "", 0)