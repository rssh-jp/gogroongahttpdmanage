package gogroongahttpdmanage

func arrayInterfaceToInt(i interface{}, indexes []int) int {
	var o interface{}
	o = i
	for i := 0; i < len(indexes); i++ {
		o = o.([]interface{})[indexes[i]]
	}
	return int(o.(float64))
}

func arrayInterfaceToFloat64(i interface{}, indexes []int) float64 {
	var o interface{}
	o = i
	for i := 0; i < len(indexes); i++ {
		o = o.([]interface{})[indexes[i]]
	}
	return o.(float64)
}

func arrayInterfaceToString(i interface{}, indexes []int) string {
	var o interface{}
	o = i
	for i := 0; i < len(indexes); i++ {
		o = o.([]interface{})[indexes[i]]
	}
	return o.(string)
}

func arrayInterfaceLen(i interface{}, indexes []int) int {
	var o interface{}
	o = i
	for i := 0; i < len(indexes); i++ {
		o = o.([]interface{})[indexes[i]]
	}
	return len(o.([]interface{}))
}

func interfaceToInt(i interface{}, index int) int {
	return int(i.([]interface{})[index].(float64))
}

func interfaceToFloat64(i interface{}, index int) float64 {
	return i.([]interface{})[index].(float64)
}

func interfaceToString(i interface{}, index int) string {
	return i.([]interface{})[index].(string)
}
