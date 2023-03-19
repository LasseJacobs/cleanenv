package cleanenv

import (
	"fmt"
)

func finalize(metaInfo []structMeta) error {
	for _, meta := range metaInfo {
		// enforce required values
		if meta.required && meta.isFieldValueZero() {
			err := fmt.Errorf("field %q is required but the value is not provided",
				meta.fieldName)
			return err
		}

		// set default values
		// ... for now this is handled by flags
		if meta.isFieldValueZero() && meta.defValue != nil {
			if err := parseValue(meta.fieldValue, *meta.defValue, meta.separator, meta.layout); err != nil {
				return err
			}
		}
	}

	return nil
}
