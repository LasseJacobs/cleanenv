package cleanenv

import (
	"flag"
	"strings"
)

func readFlagVars(metaInfo []structMeta) error {
	var values = make(map[string]*string)
	var names = make(map[string]struct{})
	var flagSet = flag.NewFlagSet("", flag.ExitOnError)
	for _, meta := range metaInfo {
		//var deflt = ""
		//if meta.defValue != nil {
		//	deflt = *meta.defValue
		//}

		flagName := toFlagName(meta.fieldName, meta.envList)
		if _, ok := names[flagName]; ok {
			// we can't have overlapping flag names
			// TODO: should we return an error here?
			continue
		}
		names[flagName] = struct{}{}
		values[meta.fieldName] = flagSet.String(flagName, "", meta.description)
	}
	flag.Parse()
	for _, meta := range metaInfo {
		flagVal := *values[meta.fieldName]
		if flagVal == "" {
			continue
		}
		if err := parseValue(meta.fieldValue, flagVal, meta.separator, meta.layout); err != nil {
			return err
		}
	}
	return nil
}

// TODO: aliases should be removed or properly supported
func toFlagName(name string, aliases []string) string {
	//TODO: the flag name needs to be improved; there should be one annotation with a name that gets reliably translated into each name (flag, env, json)
	flagName := strings.ToLower(name)
	if len(aliases) > 0 {
		flagName = strings.ReplaceAll(strings.ToLower(aliases[0]), "_", "-")
	}
	return flagName
}
