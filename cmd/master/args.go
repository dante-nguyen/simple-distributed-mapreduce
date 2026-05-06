package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	nfsRoot      = flagx.DirValue{Path: "/mnt/nfs"}
	inputPattern string
)

func prepArguments() {
	flag.IntVar(&port, "port", 8000, "grpc server port")
	flag.StringVar(&advertiseAddr, "advertise-address", "", "grpc server advertise address")
	flag.DurationVar(&healthyDuration, "healthy-duration", 30*time.Second, "maximum duration since last heartbeat of a healthy worker")
	flag.DurationVar(&healthcheckInterval, "healthcheck-interval", 5*time.Second, "interval between worker heartbeat checks")
	flag.Var(&nfsRoot, "nfs-root", "NFS root volume")
	flag.StringVar(&inputPattern, "input", "", "input files glob pattern relative to NFS root")
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
		return require("advertise address")
	case len(nfsRoot.Path) == 0:
		return require("NFS root")
	case len(inputPattern) == 0:
		return require("input pattern")
	default:
		return nil
	}
}

func require(name string) error {
	return fmt.Errorf("%s must be provided", name)
}

func exit1(err error) {
	log.Println(err.Error())
	os.Exit(1)
}

func globFiles(nfsRoot string, relPattern string) ([]string, error) {
	absPattern := filepath.Join(nfsRoot, relPattern)
	matches, err := filepath.Glob(absPattern)
	if err != nil {
		return nil, err
	}

	// will check if paths are files in server config validation
	return matches, nil
}
