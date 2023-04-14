package helper

func IsBasicDataType(str string) bool {
	basicDataTypes := []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "bool", "byte", "rune", "string"}
	for _, dataType := range basicDataTypes {
		if str == dataType {
			return true
		}
	}
	return false
}
