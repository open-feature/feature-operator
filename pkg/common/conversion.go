package common

func Bool(value bool) *bool {
	return &value
}

func Int32(value int32) *int32 {
	return &value
}

func Int64(value int64) *int64 {
	return &value
}

func String(value string) *string {
	return &value
}

func CopyString(str string) *string {
	if len(str) == 0 {
		return nil
	}
	newStr := str
	return &newStr
}

func GetBool(value *bool, defaultValue bool) *bool {
	if value != nil {
		return value
	}
	return Bool(defaultValue)
}
