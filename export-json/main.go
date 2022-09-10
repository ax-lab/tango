package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"

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

	opLoadNames := common.Start("loading names")
	names, err := data.LoadNames(importDir)
	if err != nil {
		common.ExitWithError("loading names: %v", err)
	}
	opLoadNames.Complete()
	fmt.Printf("... loaded %d names\n", len(names))

	opWriteNames := common.Start("writing names")
	outputNamesPath := path.Join(outputDir, "names.json")
	outputNamesFile, err := os.OpenFile(outputNamesPath, os.O_CREATE, os.ModePerm)
	if err != nil {
		common.ExitWithError("creating output file: %v", err)
	}
	outputNames := json.NewEncoder(outputNamesFile)
	outputNames.SetIndent("", "\t")
	if err = outputNames.Encode(names); err != nil {
		common.ExitWithError("writing output file: %v", err)
	}
	outputNamesFile.Close()
	opWriteNames.Complete()
}
