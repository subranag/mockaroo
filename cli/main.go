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
		os.Exit(2)
	}
	s := mockaroo.NewServer(conf)

	if err := s.Start(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}
}
