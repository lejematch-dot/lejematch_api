package main

import (
	"Lejematch/cmd"
	"log"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
