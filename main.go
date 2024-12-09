package main

import (
	"log"

	"github.com/nugget/roadtrip-lubelogger-sync/roadtrip"
)

func main() {
	var gt3 roadtrip.CSV

	err := gt3.LoadFile("data/CSV/2007 GT3 RS.csv")
	if err != nil {
		log.Fatal(err)
	}
}
