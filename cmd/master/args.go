package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/flagx"
)

var (
	// server
	port          int
	advertiseAddr string
	// healthcheck
	healthyDuration     time.Duration
	healthcheckInterval time.Duration
	// input
	inDir flagx.DirValue
)

func prepArguments() {
	flag.IntVar(&port, "port", 8000, "grpc server port")
	flag.StringVar(&advertiseAddr, "advertise-address", "", "grpc server advertise address")
	flag.DurationVar(&healthyDuration, "healthy-duration", 30*time.Second, "maximum duration since last heartbeat of a healthy worker")
	flag.DurationVar(&healthcheckInterval, "healthcheck-interval", 5*time.Second, "interval between worker heartbeat checks")
	flag.Var(&inDir, "in", "input directory")
	flag.Parse()

	if err := validateDirectArguments(); err != nil {
		exit1(errx.WithContext(err, "validate arguments"))
	}
}

func validateDirectArguments() error {
	switch {
	case port <= 0:
		return errors.New("invalid port")
	case len(advertiseAddr) == 0:
		return errors.New("advertise address is required")
	case len(inDir.Path) == 0:
		return errors.New("input directory is required")
	default:
		return nil
	}
}

func exit1(err error) {
	log.Println(err.Error())
	os.Exit(1)
}
