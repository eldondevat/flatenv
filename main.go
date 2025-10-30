package main

import (
	"flag"
	"fmt"
	"github.com/vincent-petithory/dataurl"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var prefix string
var dryrun bool
var perms string
var permsMode os.FileMode

var dperms string
var dpermsMode os.FileMode

func init() {
	flag.StringVar(&prefix, "prefix", "FLATENV_", "Environment variable prefix to scan for")
	flag.BoolVar(&dryrun, "dryrun", false, "Log files that would be written instead of writing them")
	flag.StringVar(&perms, "perms", "0660", "Default filesystem permissions for files")
	flag.StringVar(&dperms, "dperms", "0770", "Default filesystem permissions for directories")
}

func main() {
	flag.Parse()
	fmt.Sscanf(perms, "%o", &permsMode)
	fmt.Sscanf(dperms, "%o", &dpermsMode)
	slog.With("dperms", dpermsMode).Info("Scanned octal")
	if !strings.HasSuffix(prefix, "_") {
		slog.With("prefix", prefix).Warn("Prefix does not end with underscore, which is unexpected")
	}
	files := readEnv()
	slog.With("files", files).Info("found files")
	err := render(files)
	if err != nil {
		slog.With("error", err, "files", files).Error("Cannot render file")
		os.Exit(1)
	}
}

func readEnv() map[string]string {
	files := make(map[string]string)
	for _, environ := range os.Environ() {
		if strings.HasPrefix(environ, prefix) {
			environKV := strings.SplitN(environ, "=", 2)
			files[strings.TrimPrefix(environKV[0], prefix)] = environKV[1]
		}
	}
	return files
}

func render(files map[string]string) error {
	decodedFiles, err := translateFileNamesAndValues(files)
	if err != nil {
		return err
	}
	for fileName, contents := range decodedFiles {
		if !dryrun {
			err := os.MkdirAll(filepath.Dir(fileName), dpermsMode)
			if err != nil {
				return err
			}
			err = os.WriteFile(fileName, contents, permsMode)
			if err != nil {
				return err
			}
		} else {
			slog.With("file", fileName, "size", len(contents), "mode", fmt.Sprintf("%#o", permsMode)).Info("would render file")
		}
	}
	return nil
}

func translateFileNamesAndValues(files map[string]string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for fileSuffix, dataUri := range files {
		fileSuffixParts := strings.Split(fileSuffix, "____")
		fileName := ""
		for index, part := range fileSuffixParts {
			if index > 0 {
				fileName = fileName + "__"
			}
			fileName = fileName + strings.ReplaceAll(strings.ReplaceAll(part, "___", string(os.PathSeparator)), "__", ".")
		}
		data, err := dataurl.DecodeString(dataUri)
		if err != nil {
			return nil, err
		}
		result[fileName] = data.Data
	}
	return result, nil
}
