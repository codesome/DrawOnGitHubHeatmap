package DrawOnGitHubHeatmap

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Committer struct {
	Padding int
	Config  *Config
	i       int
}

func NewCommitter(config *Config, padding int) Committer {
	return Committer{
		Config:  config,
		Padding: padding,
		i:       0,
	}
}

func CleanRepo() error {
	if err := exec.Command("git", "checkout", "--orphan", "latest_branch").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "add", "-A").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "commit", "-am", "\"commit message\"").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "branch", "-D", "master").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "branch", "-m", "master").Run(); err != nil {
		return err
	}
	return nil
}

func (c *Committer) CommitAtIndex(i, j int) error {
	t := c.Config.IndexToTime(c.i+i, j)
	return c.CommitUsingTimestamp(t)
}

func (c *Committer) CommitUsingTimestamp(t time.Time) error {
	for i := 0; i < c.Config.EachPixelCommit; i++ {
		if err := os.Setenv("GIT_COMMITTER_DATE", t.String()); err != nil {
			return err
		}
		if err := os.Setenv("GIT_AUTHOR_DATE", t.String()); err != nil {
			return err
		}
		if err := ioutil.WriteFile("dummy.txt", []byte(t.String()+strconv.Itoa(i)), 0644); err != nil {
			return err
		}
		if err := exec.Command("git", "add", "dummy.txt").Run(); err != nil {
			fmt.Println("Place 1")
			return err
		}
		if err := exec.Command("git", "commit", "--date", t.String(), "-m", t.String()).Run(); err != nil {
			fmt.Println("Place 2")
			return err
		}
	}
	return nil
}

func (c *Committer) Next() {
	c.i += c.Padding
}

type Config struct {
	EachPixelCommit int
	Start           time.Time
}

func (c *Config) Validate() error {
	if c.EachPixelCommit <= 0 {
		return fmt.Errorf("Pixel commits should be >0")
	}
	if c.Start.Weekday() != time.Sunday {
		return fmt.Errorf("Start time should be Sunday")
	}
	return nil
}

func (c *Config) IndexToTime(i, j int) time.Time {
	days := (i * Height) + j
	return c.Start.AddDate(0, 0, days)
}
