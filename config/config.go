package config

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config holds the configuration values for an anime.
type Config struct {
	Name          string
	VideoRegex    *regexp.Regexp
	SubtitleRegex *regexp.Regexp
}

func (c *Config) String() string {
	return fmt.Sprintf("(%s)\n\tVideoRegex: %v\n\tSubtitleRegex: %v", c.Name, c.VideoRegex, c.SubtitleRegex)
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	data := map[string]string{}
	if err := value.Decode(data); err != nil {
		return errors.WithStack(err)
	}
	c.Name = data["name"]
	c.VideoRegex = regexp.MustCompile(data["videoRegex"])
	c.SubtitleRegex = regexp.MustCompile(data["subtitleRegex"])
	return nil
}

// LoadConfig loads the configuration file for an anime.
func LoadConfig(filename string) ([]*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cfg := make([]*Config, 0)
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return nil, errors.WithStack(err)
	}
	return cfg, nil
}
