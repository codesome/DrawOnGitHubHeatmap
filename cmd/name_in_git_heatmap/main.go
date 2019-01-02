package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	nameingitheatmap "github.com/codesome/name_in_git_heatmap"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "The Prometheus monitoring server")

	a.HelpFlag.Short('h')

	var (
		pixelConfigFile string
		mainString      string
		commitsPerDay   int
	)

	a.Flag("pixel-layout-config", "Configuration file for pixel layout.").
		Default("pixelConfig.yaml").StringVar(&pixelConfigFile)

	a.Flag("commits-per-day", "Configuration file for pixel layout.").
		Default("5").IntVar(&commitsPerDay)

	a.Flag("text", "Text to write in the heatmap.").
		Default("").StringVar(&mainString)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if err := nameingitheatmap.CleanRepo(); err != nil {
		panic(err)
	}

	pixelConfig, err := nameingitheatmap.ParsePixelConfigFromFile(pixelConfigFile)
	if err != nil {
		panic(err)
	}

	padding := 0
	if len(mainString)*pixelConfig.Width > nameingitheatmap.Width {
		panic(fmt.Errorf("Too many characters"))
	}
	if len(mainString)*(pixelConfig.Width+1) <= nameingitheatmap.Width {
		padding = pixelConfig.Width + 1
	}
	if padding == 0 {
		padding = pixelConfig.Width
	}

	initialOffset := (nameingitheatmap.Width - len(mainString)*(pixelConfig.Width+1) + 1) / 2

	now := time.Now()
	config := &nameingitheatmap.Config{
		EachPixelCommit: commitsPerDay,
	}
	config.Start = now.AddDate(0, 0, -int(now.Weekday()))
	config.Start = config.Start.AddDate(0, 0, (initialOffset-nameingitheatmap.Width)*nameingitheatmap.Height)

	if err := config.Validate(); err != nil {
		panic(err)
	}

	committer := nameingitheatmap.NewCommitter(config, padding)

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
