package nameingitheatmap

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

const (
	EmptyPixel  rune = '.'
	FilledPixel rune = '*'
	Height      int  = 7
	Width       int  = 52
)

func ParsePixelConfigFromFile(filename string) (*PixelConfig, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParsePixelConfigFromBytes(b)
}

func ParsePixelConfigFromBytes(b []byte) (*PixelConfig, error) {
	var rawPixelConfig RawPixelConfig
	if err := yaml.UnmarshalStrict(b, &rawPixelConfig); err != nil {
		return nil, err
	}

	pixelConfig := &PixelConfig{
		Width:      rawPixelConfig.Width,
		Characters: make(map[rune][][2]int, len(rawPixelConfig.Characters)),
	}

	for _, rawChar := range rawPixelConfig.Characters {
		if len(rawChar.Layout) != Height {
			return nil, fmt.Errorf("Layout height should be %d, got %d. Char:%c", Height, len(rawChar.Layout), rawChar.Char)
		}
		for j, line := range rawChar.Layout {
			if len(line) != pixelConfig.Width {
				return nil, fmt.Errorf("Layout width doesn't match config width. Char:%c, LineIndex:%d", rawChar.Char, j)
			}
			for i, pixel := range line {
				if rune(pixel) == FilledPixel {
					pixelConfig.Characters[rune(rawChar.Char[0])] = append(pixelConfig.Characters[rune(rawChar.Char[0])], [2]int{i, j})
				}
			}
		}
	}

	return pixelConfig, nil
}

type PixelConfig struct {
	Width      int
	Characters map[rune][][2]int
}

type RawPixelConfig struct {
	Width      int             `yaml:"width"`
	Characters []RawCharLayout `yaml:"characters"`
}

type RawCharLayout struct {
	Char   string   `yaml:"char"`
	Layout []string `yaml:"layout"`
}
