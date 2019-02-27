package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	au "github.com/logrusorgru/aurora"
)

// delete output files in ${outputdir} older than ${retain} days
func purgeOutput(outputdir string, retain int) error {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.json", outputdir))
	if err != nil {
		return fmt.Errorf("can't glob output files (%s)", err.Error())
	}

	now := time.Now().Unix()
	for _, file := range files {
		basename := strings.TrimSuffix(file, ".json")
		basename = strings.TrimPrefix(basename, fmt.Sprintf("%s/", outputdir))
		i, err := strconv.ParseInt(basename, 0, 64)
		if err != nil {
			fmt.Printf("can't obtain file creation time for %s (can't parse %s)\n", file, basename)
			continue
		}
		if (int(now) - int(i)) > (retain * 86400) {
			err = os.Remove(file)
			if err != nil {
				fmt.Printf("can't delete file %s\n", file)
				continue
			}
			fmt.Printf("%s old output file %s\n", au.Red(au.Bold("Removed")), file)
		}
	}

	return nil
}
