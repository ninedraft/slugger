package main

import (
	"math"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

type gitStats struct {
	CommitStats []int
	Branch      string
}

func getCommits(wd string) gitStats {
	var stats gitStats
	var repo, errOpenRepo = searchAndOpenRepo(wd)
	if errOpenRepo != nil {
		return stats
	}
	var head, errGetHead = repo.Head()
	if errGetHead == nil && head.Name().IsBranch() {
		stats.Branch = head.Name().Short()
	}
	// log sorted by commit timestamp!
	var gitLog, errOpenLog = repo.Log(&git.LogOptions{
		Order: git.LogOrderDefault,
	})
	if errOpenLog != nil {
		return stats
	}
	const day = 24 * time.Hour
	const historyDeptthInDays = 7
	var now = time.Now().Round(day)
	var fromDate = now.Add(-historyDeptthInDays * day).Round(day)
	var commitStats = make([]int, historyDeptthInDays+1)
	for {
		var commit, errNextCommit = gitLog.Next()
		if errNextCommit != nil {
			break
		}
		var commitDate = commit.Author.When
		if commitDate.Before(fromDate) {
			break
		}
		var day = int(commitDate.Sub(fromDate) / day)
		commitStats[day]++
	}
	stats.CommitStats = commitStats
	return stats
}

func searchAndOpenRepo(wd string) (*git.Repository, error) {
	wd = strings.TrimSuffix(wd, string(filepath.Separator))
	var repo, errOpen = git.PlainOpen(wd)
	switch errOpen {
	case git.ErrRepositoryNotExists:
		var parentDir, _ = filepath.Split(wd)
		if parentDir != "" && parentDir != "/" {
			return searchAndOpenRepo(parentDir)
		}
		return nil, errOpen
	case nil:
		return repo, nil
	default:
		return nil, errOpen
	}
}

func levelToColor(level int) []ASCIICode {
	switch level {
	case 1, 2:
		return []ASCIICode{Red}
	case 3, 4, 5:
		return []ASCIICode{Green}
	case 6, 7:
		return []ASCIICode{Blue}
	default:
		return []ASCIICode{}
	}
}

func levelToMarker(level int) string {
	var markers = []string{" ", "_", "▃", "▄", "▅", "▆", "▇"}
	var N = len(markers)
	if level >= N {
		level = N - 1
	}
	if level < 0 {
		level = 0
	}
	return markers[level]
}

func commitStatsToLevelMarkers(commitStats []int) []string {
	var markers []string
	for _, level := range commitStatsToLevels(commitStats) {
		markers = append(markers, levelToMarker(level))
	}
	return markers
}

func commitStatsToLevels(commitStats []int) []int {
	var maxCommits = float64(max(commitStats))
	const maxLevel = 7
	var levels []int
	for _, stat := range commitStats {
		var level = int(math.Round(maxLevel * float64(stat) / maxCommits))
		if level == 0 && stat > 0 {
			level = 1
		}
		levels = append(levels, level)
	}
	return levels
}

func sum(commitStats []int) int {
	var sum int
	for _, stat := range commitStats {
		sum += stat
	}
	return sum
}

func max(commitStats []int) int {
	var N = len(commitStats)
	switch N {
	case 0:
		return 0
	case 1:
		return commitStats[0]
	}
	var max = commitStats[0]
	for _, stat := range commitStats[1:] {
		if max < stat {
			max = stat
		}
	}
	return max
}
