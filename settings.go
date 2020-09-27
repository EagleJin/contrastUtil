package main

import "flag"

// AppSettings is the struct of main configuration
type AppSettings struct {
	filePath     string
	diffResult   string
	ignoreFields string
}

var Settings AppSettings

func init() {
	flag.StringVar(&Settings.filePath, "file-path", "", "replay result file path. Example: `/tmp/replay_result.log`")
	flag.StringVar(&Settings.diffResult, "output-file", "", "compare result file path. Example: `/tmp/diff.log`")
	flag.StringVar(&Settings.ignoreFields, "ignore-fields", "", "Specify the fields to ignore in the results.")
}
