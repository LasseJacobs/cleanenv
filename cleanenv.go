package cleanenv

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	// DefaultSeparator is a default list and map separator character
	DefaultSeparator = ","
)

// Supported tags
const (
	// TagEnv Name of the environment variable or a list of names
	TagEnv = "env"
	// TagEnvLayout Value parsing layout (for types like time.Time)
	TagEnvLayout = "env-layout"
	// TagEnvDefault Default value
	TagEnvDefault = "env-default"
	// TagEnvSeparator Custom list and map separator
	TagEnvSeparator = "env-separator"
	// TagEnvDescription Environment variable description
	TagEnvDescription = "env-description"
	// TagEnvRequired Flag to mark a field as required
	TagEnvRequired = "env-required"
	// TagEnvPrefix Flag to specify prefix for structure fields
	TagEnvPrefix = "env-prefix"
)

// Setter is an interface for a custom value setter.
//
// To implement a custom value setter you need to add a SetValue function to your type that will receive a string raw value:
//
//	type MyField string
//
//	func (f *MyField) SetValue(s string) error {
//		if s == "" {
//			return fmt.Errorf("field value can't be empty")
//		}
//		*f = MyField("my field is: " + s)
//		return nil
//	}
type Setter interface {
	SetValue(string) error
}

// ReadConfig reads configuration file and parses it depending on tags in structure provided.
// Then it reads and parses
//
// Example:
//
//	type ConfigDatabase struct {
//		Port     string `yaml:"port" env:"PORT" env-default:"5432"`
//		Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
//		Name     string `yaml:"name" env:"NAME" env-default:"postgres"`
//		User     string `yaml:"user" env:"USER" env-default:"user"`
//		Password string `yaml:"password" env:"PASSWORD"`
//	}
//
//	var cfg ConfigDatabase
//
//	err := cleanenv.ReadConfig("config.yml", "PREFIX", &cfg)
//	if err != nil {
//	    ...
//	}
func ReadConfig(path string, appname string, cfg interface{}) error {
	var err error
	if path != "" {
		err = parseFile(path, cfg)
		if err != nil {
			return err
		}
	}

	var meta []structMeta
	meta, err = readStructMetadata(cfg)
	if err != nil {
		return err
	}

	err = readEnvVars(meta, toPrefix(appname))
	if err != nil {
		return err
	}

	err = readFlagVars(meta)
	if err != nil {
		return err
	}

	return finalize(meta)
}

// ReadEnv reads environment variables into the structure.
func ReadEnv(cfg interface{}, appname string) error {
	meta, err := readStructMetadata(cfg)
	if err != nil {
		return err
	}

	err = readEnvVars(meta, toPrefix(appname))
	if err != nil {
		return err
	}

	return finalize(meta)
}

// GetDescription returns a description of environment variables.
// You can provide a custom header text.
func GetDescription(cfg interface{}, headerText *string) (string, error) {
	meta, err := readStructMetadata(cfg)
	if err != nil {
		return "", err
	}

	var header, description string

	if headerText != nil {
		header = *headerText
	} else {
		header = "Environment variables:"
	}

	for _, m := range meta {
		if len(m.envList) == 0 {
			continue
		}

		for idx, env := range m.envList {

			elemDescription := fmt.Sprintf("\n  %s %s", env, m.fieldValue.Kind())
			if idx > 0 {
				elemDescription += fmt.Sprintf(" (alternative to %s)", m.envList[0])
			}
			elemDescription += fmt.Sprintf("\n    \t%s", m.description)
			if m.defValue != nil {
				elemDescription += fmt.Sprintf(" (default %q)", *m.defValue)
			}
			description += elemDescription
		}
	}

	if description != "" {
		return header + description, nil
	}
	return "", nil
}

// Usage returns a configuration usage help.
// Other usage instructions can be wrapped in and executed before this usage function.
// The default output is STDERR.
func Usage(cfg interface{}, headerText *string, usageFuncs ...func()) func() {
	return FUsage(os.Stderr, cfg, headerText, usageFuncs...)
}

// FUsage prints configuration help into the custom output.
// Other usage instructions can be wrapped in and executed before this usage function
func FUsage(w io.Writer, cfg interface{}, headerText *string, usageFuncs ...func()) func() {
	return func() {
		for _, fn := range usageFuncs {
			fn()
		}

		_ = flag.Usage

		text, err := GetDescription(cfg, headerText)
		if err != nil {
			return
		}
		if len(usageFuncs) > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w, text)
	}
}
