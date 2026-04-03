package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nlduy0310/simple-distributed-mapreduce/errorsx"
	"github.com/nlduy0310/simple-distributed-mapreduce/worker"
	"github.com/nlduy0310/simplelog"
)

var logger = simplelog.NewLogger("entrypoint", simplelog.DEBUG)

func main() {
	os.Exit(run())
}

func run() int {
	if err := godotenv.Load(".env"); err != nil {
		logger.Fatal(errorsx.WrapAsMessage("failed to load env file", err))
	}

	svr, err := worker.Setup()
	if err != nil {
		logger.Fatal(errorsx.WrapAsMessage("can not setup worker", err))
	}
	defer svr.Close()

	if err = svr.Serve(); err != nil {
		logger.Error(errorsx.WrapAsMessage("server stopped with error", err))
		return 1
	}

	return 0
}
