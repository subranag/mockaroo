package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/subranag/mockaroo"
)

func main() {
	mockConfig := flag.String("conf", "", "the mockaroo config file")
	flag.Parse()

	if len(*mockConfig) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// parse config
	conf, err := mockaroo.LoadConfig(mockConfig)
	if err != nil {
		log.Fatalf("error loading config :%v", err)
		os.Exit(2)
	}
	s := mockaroo.NewServer(conf)

	if err := s.Start(); err != nil {
		log.Fatalf("error starting server :%v", err)
		os.Exit(2)
	}
}
