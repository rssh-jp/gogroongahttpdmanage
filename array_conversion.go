package gogroongahttpdmanage

func arrayInterfaceToInt(i interface{}, index int) int {
	return int(i.([]interface{})[index].(float64))
}

func arrayInterfaceToFloat64(i interface{}, index int) float64 {
	return i.([]interface{})[index].(float64)
}

func arrayInterfaceToString(i interface{}, index int) string {
	return i.([]interface{})[index].(string)
}

