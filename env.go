package cleanenv

import (
	"os"
)

// readEnvVars reads environment variables to the provided configuration structure
func readEnvVars(metaInfo []structMeta, prefix string) error {
	for _, meta := range metaInfo {
		var rawValue *string

		//TODO: env list needs to become a single value; what is the use-case for multi-value
		for _, env := range meta.envList {
			if value, ok := os.LookupEnv(prefix + env); ok {
				rawValue = &value
				break
			}
		}

		if rawValue == nil {
			continue
		}

		if err := parseValue(meta.fieldValue, *rawValue, meta.separator, meta.layout); err != nil {
			return err
		}
	}

	return nil
}
