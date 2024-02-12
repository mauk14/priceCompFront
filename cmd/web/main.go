package main

import "log"

func main() {
	err := Route().Run(":5000")
	if err != nil {
		log.Fatal(err)
	}
}
