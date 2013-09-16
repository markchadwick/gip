package main

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	if err := get(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
