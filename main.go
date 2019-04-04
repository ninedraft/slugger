package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"text/template"
	"time"

	"github.com/mitchellh/go-ps"
)

func main() {
	defer func() {
		var err = recover()
		if err != nil {
			debug.PrintStack()
			fmt.Printf("%v > ", err)
		}
	}()
	var now = time.Now()
	var timestamp = ColoredText(now.Format(time.RFC822), Yellow, Bold)
	var user = ColoredText(os.Getenv("USER"), Blue, Bold)
	var workingDirPath, _ = os.Getwd()
	var repoStats = getCommits(workingDirPath)
	var workingDirPathWithMagic = strings.Replace(workingDirPath, os.Getenv("HOME"), "~", 1)
	var coloredWorkingDir = ColoredText(workingDirPathWithMagic, Green)

	var commitLine = &bytes.Buffer{}
	if len(repoStats.CommitStats) > 0 {
		commitLine.WriteString("git:|")
		for _, level := range commitStatsToLevels(repoStats.CommitStats) {
			var marker = levelToMarker(level)
			var color = levelToColor(level)
			commitLine.WriteString(ColoredText(marker, color...))
		}
		commitLine.WriteString("|")
	}
	var branchMarker string
	switch repoStats.Branch {
	case "master":
		branchMarker = ColoredText(repoStats.Branch, Red)
	case "":
		// nothing
	default:
		branchMarker = ColoredText(repoStats.Branch, Green, Bold)
	}
	var parentPID = os.Getppid()
	var parentProccess, errGetProcess = ps.FindProcess(parentPID)
	var parentName string
	if errGetProcess == nil {
		parentName = parentProccess.Executable()
	}
	var fstr = "{{.Username}} {{.Timestamp}} {{.Wd}} {{.CommitStats}} {{.Branch}}\n{{.Shell}}> "
	format(os.Stdout, fstr, map[string]interface{}{
		"Username":    user,
		"Timestamp":   timestamp,
		"Wd":          coloredWorkingDir,
		"CommitStats": commitLine.String(),
		"Branch":      branchMarker,
		"Shell":       parentName,
	})
}

func format(wr io.Writer, text string, data interface{}) {
	var prompt = template.New("prompt")
	if err := template.Must(prompt.Parse(text)).
		Execute(wr, data); err != nil {
		panic(err)
	}
}
