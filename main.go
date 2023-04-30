package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/chuckha/naming/config"
	"github.com/pkg/errors"
)

func main() {
	directory := flag.String("dir", "", "directory to use")
	dry := flag.Bool("dryrun", true, "set dryrun to false to execute rename")
	flag.Parse()

	if *directory == "" {
		fmt.Println("pass in a directory to scan")
		os.Exit(0)
	}
	dryrun := *dry
	dir := *directory
	c, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	for _, cfg := range c {
		eps, err := FindEpisodes(dir, cfg)
		if err != nil {
			panic(err)
		}
		for _, e := range eps {
			newVideo := path.Join(path.Dir(e.VideoFile), e.VideoFilename())
			newSubtitle := path.Join(path.Dir(e.SubtitleFile), e.SubtitleFilename())

			if dryrun {
				fmt.Printf("would move %q to %q\n", e.VideoFile, newVideo)
				fmt.Printf("would move %q to %q\n", e.SubtitleFile, newSubtitle)
				continue
			}
			if e.VideoExt != "" {
				if err := renameFile(e.VideoFile, newVideo); err != nil {
					panic(err)
				}
			}
			if e.SubtitleExt != "" {
				if err := renameFile(e.SubtitleFile, newSubtitle); err != nil {
					panic(err)
				}
			}
		}
	}
}

func renameFile(oldPath, newPath string) error {
	_, err := os.Stat(newPath)
	if err == nil {
		fmt.Printf("skipping existing file %s\n", newPath)
		return nil
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Episode is a single show episode with a matching subtitle file.
type Episode struct {
	Name         string
	VideoFile    string
	SubtitleFile string
	Season       int
	Episode      int
	VideoExt     string
	SubtitleExt  string
}

func (e *Episode) Basename() string {
	// TODO: consider customizing output
	return fmt.Sprintf("%sS%02dE%02d", e.Name, e.Season, e.Episode)
}

func (e *Episode) VideoFilename() string {
	return fmt.Sprintf("%s.%s", e.Basename(), e.VideoExt)
}
func (e *Episode) SubtitleFilename() string {
	return fmt.Sprintf("%s.%s", e.Basename(), e.SubtitleExt)
}

type filenameData struct {
	season, episode int
	ext             string
}

type key struct {
	season, episode int
}

// FindFiles finds all video files and their corresponding subtitle files in a directory using the provided configuration.
func FindEpisodes(dir string, config *config.Config) (map[key]*Episode, error) {
	episodes := map[key]*Episode{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if config.VideoRegex.MatchString(info.Name()) {
			// create or get episode
			match := config.VideoRegex.FindStringSubmatch(info.Name())
			fd := filenameData{}
			for i, name := range config.VideoRegex.SubexpNames() {
				switch name {
				case "episode":
					i, _ := strconv.Atoi(match[i])
					fd.episode = i
				case "season":
					i, _ := strconv.Atoi(match[i])
					fd.season = i
				case "ext":
					fd.ext = match[i]
				case "":
				default:
					panic(fmt.Sprintf("unknown named regular expression %v", name))
				}
			}
			// Default to season 1 if there is no season regex defined
			if fd.season == 0 {
				fd.season = 1
			}
			k := key{season: fd.season, episode: fd.episode}
			if _, ok := episodes[k]; !ok {
				episodes[k] = &Episode{}
			}
			episodes[k].Name = config.Name
			episodes[k].VideoFile = path
			episodes[k].Season = fd.season
			episodes[k].Episode = fd.episode
			episodes[k].VideoExt = fd.ext
			// should not match both subtitle file and video file
			return nil
		}

		if config.SubtitleRegex.MatchString(info.Name()) {
			// create or get episode
			match := config.SubtitleRegex.FindStringSubmatch(info.Name())
			fd := filenameData{}
			for i, name := range config.SubtitleRegex.SubexpNames() {
				switch name {
				case "episode":
					i, _ := strconv.Atoi(match[i])
					fd.episode = i
				case "season":
					i, _ := strconv.Atoi(match[i])
					fd.season = i
				case "ext":
					fd.ext = match[i]
				case "":
				default:
					panic(fmt.Sprintf("unknown named regular expression %v", name))
				}
			}
			// Default to season 1 if there is no season regex defined
			if fd.season == 0 {
				fd.season = 1
			}
			k := key{season: fd.season, episode: fd.episode}
			if _, ok := episodes[k]; !ok {
				episodes[k] = &Episode{}
			}
			episodes[k].Name = config.Name
			episodes[k].SubtitleFile = path
			episodes[k].Season = fd.season
			episodes[k].Episode = fd.episode
			episodes[k].SubtitleExt = fd.ext
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return episodes, nil
}
