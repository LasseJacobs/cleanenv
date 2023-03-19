package cleanenv

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// parseFile parses configuration file according to it's extension
//
// Currently following file extensions are supported:
//
// - yaml
//
// - json
func parseFile(path string, cfg interface{}) error {
	// open the configuration file
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	// parse the file depending on the file type
	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".yaml", ".yml":
		err = parseYAML(f, cfg)
	case ".json":
		err = parseJSON(f, cfg)
	default:
		return fmt.Errorf("file format '%s' doesn't supported by the parser", ext)
	}
	if err != nil {
		return fmt.Errorf("config file parsing error: %s", err.Error())
	}
	return nil
}

// parseYAML parses YAML from reader to data structure
func parseYAML(r io.Reader, str interface{}) error {
	return yaml.NewDecoder(r).Decode(str)
}

// parseJSON parses JSON from reader to data structure
func parseJSON(r io.Reader, str interface{}) error {
	return json.NewDecoder(r).Decode(str)
}
