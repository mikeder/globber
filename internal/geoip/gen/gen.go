package main

import (
	"log"
	"os"

	"github.com/mikeder/globber/internal/geoip"
)

func main() {
	if err := geoip.WriteRecordsGob("data.gob"); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
