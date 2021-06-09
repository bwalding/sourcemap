package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
)

type SourceMap struct {
	FileNames []string `json:"sources"`
	Contents  []string `json:"sourcesContent"`
}

type Parser struct {
	Maps   chan RawMap
	Logger *zap.Logger
	Output string
}

func (p *Parser) Run() {
	for raw := range p.Maps {
		err := p.parse(raw)
		if err != nil {
			p.Logger.Error("cannot parse source map", zap.Error(err))
			continue
		}
	}
}

func (p *Parser) parse(raw RawMap) error {
	var m SourceMap
	err := json.Unmarshal(raw.Content, &m)
	if err != nil {
		return fmt.Errorf("read JSON: %v", err)
	}
	for i, fname := range m.FileNames {
		fname = strings.ReplaceAll(fname, "../", ".")
		fname = strings.ReplaceAll(fname, "webpack://", "")
		fname = strings.ReplaceAll(fname, "://", "")
		fname = path.Join(p.Output, raw.Host, fname)

		if i >= len(m.Contents) {
			return errors.New("sources is longer than sourcesContent")
		}
		if strings.HasPrefix(fname, "external ") {
			return errors.New("external source maps unsupported")
		}

		parent, _ := path.Split(fname)
		err = os.MkdirAll(parent, 0770)
		if err != nil {
			return fmt.Errorf("create dir: %v", err)
		}

		p.Logger.Debug("writing file", zap.String("path", fname))
		err = os.WriteFile(fname, []byte(m.Contents[i]), 0660)
		if err != nil {
			return fmt.Errorf("write file: %v", err)
		}
	}
	return nil
}
