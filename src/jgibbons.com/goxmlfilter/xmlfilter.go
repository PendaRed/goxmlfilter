package main

import (
	"fmt"
	"jgibbons.com/goxmlfilter/config"
	"jgibbons.com/goxmlfilter/filterservice"
	"os"
)

func main() {
	//	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) != 1 {
		usage()
	}

	config, err := config.New(argsWithoutProg[0])
	if err != nil {
		fmt.Println("Failed to read config, aborting")
		usage()
	}
	config.DumpConfig()
	err = filterservice.Start(config)
	if err != nil {
		fmt.Printf("Failed to process input data, %v\n", err)
		usage()
	}
}

func usage() {
	message := `Usage xmlfilter <config file.json>
  where config file.json holds the config as shown in the example.json`
	fmt.Printf("%s", message)

	os.Exit(-1)
}
