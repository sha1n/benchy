package cli

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const (
	// ArgNameConfig : program arg name
	ArgNameConfig = "config"
	// ArgNameOutputFile : program arg name
	ArgNameOutputFile = "out-file"
	// ArgNameFormat : program arg name
	ArgNameFormat = "format"
	// ArgNamePipeStdout : program arg name
	ArgNamePipeStdout = "pipe-stdout"
	// ArgNamePipeStderr : program arg name
	ArgNamePipeStderr = "pipe-stderr"
	// ArgNameDebug : program arg name
	ArgNameDebug = "debug"
	// ArgNameSilent : program arg name
	ArgNameSilent = "silent"
	// ArgNameLabel : program arg name
	ArgNameLabel = "label"
	// ArgNameHeaders : program arg name
	ArgNameHeaders = "headers"

	// ArgValueReportFormatTxt : Plain text report format arg value
	ArgValueReportFormatTxt = "txt"
	// ArgValueReportFormatCsv : CSV report format arg value
	ArgValueReportFormatCsv = "csv"
	// ArgValueReportFormatCsvRaw : CSV raw data report format value
	ArgValueReportFormatCsvRaw = "csv/raw"
	// ArgValueReportFormatMarkdown : Markdown report format arg value
	ArgValueReportFormatMarkdown = "md"
	// ArgValueReportFormatMarkdownRaw : Markdown report format arg value
	ArgValueReportFormatMarkdownRaw = "md/raw"
)

// ResolveOutputFileArg resolves an output file argument based on user input.
// If the specified argument is empty, stdout is returned.
func ResolveOutputFileArg(cmd *cobra.Command, name string) *os.File {
	var outputFile = os.Stdout
	var err error = nil

	if outputFilePath := GetString(cmd, name); outputFilePath != "" {
		resolvedfilePath := expandPath(outputFilePath)
		outputFile, err = os.OpenFile(resolvedfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	CheckInitFatal(err)

	return outputFile
}

func GetString(cmd *cobra.Command, name string) string {
	v, err := cmd.Flags().GetString(name)
	CheckArgFatal(err)

	return v
}

func GetBool(cmd *cobra.Command, name string) bool {
	v, err := cmd.Flags().GetBool(name)
	CheckArgFatal(err)

	return v
}

func GetStringSlice(cmd *cobra.Command, name string) []string {
	v, err := cmd.Flags().GetStringSlice(ArgNameLabel)
	CheckArgFatal(err)

	return v
}

// TODO this has been copied from pgk/command_exec.go. Maybe share or use an existing implementation if exists.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		if p, err := os.UserHomeDir(); err == nil {
			return filepath.Join(p, path[1:])
		}
		log.Warnf("Failed to resolve user home for path '%s'", path)
	}

	return path
}