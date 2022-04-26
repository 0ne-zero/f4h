package setting

import (
	"errors"

	"github.com/0ne-zero/f4h/constansts"
)

func ReadFieldsInSettingFile(fields_name []string) (map[string]string, error) {
	var fields_value map[string]string
	for _, fn := range fields_name {
		value, exists := constansts.SettingData[fn]
		if !exists {
			return nil, errors.New("Key doesn't exists in setting data")
		}
		fields_value[fn] = value
	}
	return fields_value, nil
}
func ReadFieldInSettingData(field_name string) (string, error) {
	value, exists := constansts.SettingData[field_name]
	if !exists {
		return "", errors.New("Key doesn't exists in setting data")
	}
	return value, nil
}
