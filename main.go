package main

import (
	"fmt"
	"log"
	"os"
	fp "wget/packages/flag-parser"
)

func main() {
	// example of usage package
	parser := fp.CreateParser()

	parser.Add("B", "to use program in backgound", true)
	parser.Add("rate-limit", "to limit download speed", false)

	storage, err := parser.Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	hasWrongFlag := storage.HasFlag("A")
	hasRightFlag := storage.HasFlag("B")

	fmt.Printf("Is wrong flag found: %v\nIs right flag found: %v\n", hasWrongFlag, hasRightFlag)

	flag, err := storage.GetFlag("rate-limit")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("value of rate-limit flag: %s\n", flag.GetValue())
}
