package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/subranag/mockaroo"
)

func main() {
	mockConfig := flag.String("conf", "", "the mockaroo config file")
	flag.Parse()

	if len(*mockConfig) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	fmt.Printf("starting server with config file : %s\n", *mockConfig)
	// parse config
	conf, err := mockaroo.LoadConfig(mockConfig)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	s := mockaroo.NewServer(conf)
	s.Start()
}
