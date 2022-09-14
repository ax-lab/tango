package main

import (
	"flag"
	"fmt"

	"github.com/ax-lab/tango/common"
	"github.com/ax-lab/tango/export-json/data"
)

func main() {
	importDirPtr := flag.String("import", "", "database directory")
	outputDir := common.GetOutputDir("output", "output directory for the JSON data")
	importDir := *importDirPtr

	if importDir == "" {
		common.ExitWithError("invalid arguments: missing import directory")
	}

	fmt.Printf("Exporting JSON data from `%s` to `%s`...\n", importDir, outputDir)

	var err error

	opLoadNames := common.Start("exporting names")
	err = data.ExportNames(importDir, outputDir)
	if err != nil {
		common.ExitWithError("exporting names: %v", err)
	}
	opLoadNames.Complete()

	opLoadEntries := common.Start("exporting entries")
	err = data.ExportEntries(importDir, outputDir)
	if err != nil {
		common.ExitWithError("exporting entries: %v", err)
	}
	opLoadEntries.Complete()
}
