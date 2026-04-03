package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nlduy0310/simple-distributed-mapreduce/cli"
	"github.com/nlduy0310/simple-distributed-mapreduce/errorsx"
	"github.com/nlduy0310/simple-distributed-mapreduce/master"
	"github.com/nlduy0310/simplelog"
)

var logger = simplelog.NewLogger("entrypoint", simplelog.DEBUG)

func main() {
	os.Exit(run())
}

func run() int {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal(errorsx.WrapAsMessage("failed to load env file", err))
	}

	opts, err := cli.ParseMasterCLIOptions()
	if err != nil {
		logger.Fatal(errorsx.WrapAsMessage("failed to parse CLI options", err))
	}

	svr, err := master.Setup(opts)
	if err != nil {
		logger.Fatal(errorsx.WrapAsMessage("failed to setup server", err))
	}

	if err := svr.Serve(); err != nil {
		logger.Error(errorsx.WrapAsMessage("server stopped with error", err))
		return 1
	}

	return 0
}
