package main

import (
	"command-line/app"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Start")

	application := app.Create()
	if error := application.Run(os.Args); error != nil {
		log.Fatal(error)
	}

}
