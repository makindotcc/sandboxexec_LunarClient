package main

import (
	"log"

	"github.com/makindotcc/lunarlauncher"
)

func main() {
	log.Println("Preparing lunar assets...")
	log.Println("Fetching launch meta...")
	launchMeta, err := lunarlauncher.FetchLaunchMeta(lunarlauncher.Mc1_16, "master")
	if err != nil {
		panic(err)
	}
	log.Println("Got launch meta:", launchMeta)
	
}
