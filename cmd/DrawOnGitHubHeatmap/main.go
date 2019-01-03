package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/codesome/DrawOnGitHubHeatmap"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "Draw on GitHub heatmap")

	a.HelpFlag.Short('h')

	var (
		pixelConfigFile string
		mainString      string
		commitsPerDay   int
	)

	a.Flag("pixel-layout-config", "Configuration file for pixel layout.").
		Default("pixelConfig.yaml").StringVar(&pixelConfigFile)

	a.Flag("commits-per-day", "Number of commits on each day on heatmap. This should be >= the number of commits on darkest spot of your heatmap.").
		Default("5").IntVar(&commitsPerDay)

	a.Flag("text", "Text to write in the heatmap.").
		Default("").StringVar(&mainString)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if len(mainString) == 0 {
		return
	}

	if err := DrawOnGitHubHeatmap.CleanRepo(); err != nil {
		panic(err)
	}

	pixelConfig, err := DrawOnGitHubHeatmap.ParsePixelConfigFromFile(pixelConfigFile)
	if err != nil {
		panic(err)
	}

	padding := 0
	if len(mainString)*pixelConfig.Width > DrawOnGitHubHeatmap.Width {
		panic(fmt.Errorf("Too many characters"))
	}
	if len(mainString)*(pixelConfig.Width+1) <= DrawOnGitHubHeatmap.Width {
		padding = pixelConfig.Width + 1
	}
	if padding == 0 {
		padding = pixelConfig.Width
	}

	initialOffset := (DrawOnGitHubHeatmap.Width - len(mainString)*(pixelConfig.Width+1) + 1) / 2

	now := time.Now()
	config := &DrawOnGitHubHeatmap.Config{
		EachPixelCommit: commitsPerDay,
	}
	config.Start = now.AddDate(0, 0, -int(now.Weekday()))
	config.Start = config.Start.AddDate(0, 0, (initialOffset-DrawOnGitHubHeatmap.Width)*DrawOnGitHubHeatmap.Height)

	if err := config.Validate(); err != nil {
		panic(err)
	}

	committer := DrawOnGitHubHeatmap.NewCommitter(config, padding)

	for _, c := range mainString {
		indices := pixelConfig.Characters[rune(c)]
		for _, idx := range indices {
			if err := committer.CommitAtIndex(idx[0], idx[1]); err != nil {
				panic(err)
			}
		}
		committer.Next()
	}
}
