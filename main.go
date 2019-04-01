package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
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
	fmt.Printf("%s %s %s %s %s\n> ", user, timestamp, coloredWorkingDir, commitLine, branchMarker)
}
