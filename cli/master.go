package cli

import "github.com/alecthomas/kingpin/v2"

type MasterCLIOptions struct {
	FilePaths []string
}

func ParseMasterCLIOptions() (MasterCLIOptions, error) {
	// gotta find a better way
	paths := kingpin.Arg("inputs", "input files").Required().Strings()
	kingpin.Parse()

	ret := MasterCLIOptions{
		FilePaths: *paths,
	}
	return ret, nil
}
