package gogroongahttpdmanage

func parseHeader(o interface{}) (h Header) {
	h.ReturnCode = arrayInterfaceToInt(o, []int{0})
	h.UnixTimeWhenCommandIsStarted = arrayInterfaceToFloat64(o, []int{1})
	h.ElapsedTime = arrayInterfaceToFloat64(o, []int{2})

	if h.ReturnCode >= 0 {
		return
	}

	if arrayInterfaceLen(o, []int{}) < 3 {
		return
	}

	h.ErrorMessage = arrayInterfaceToString(o, []int{3})

	if arrayInterfaceLen(o, []int{}) < 4 {
		return
	}

	h.ErrorLocation = parseErrorLocation(o.([]interface{})[4])

	return
}

func parseErrorLocation(o interface{}) (el ErrorLocation) {
	if arrayInterfaceLen(o, []int{}) < 1 {
		return
	}

	el.LocationInGroonga = parseLocationInGroonga(o.([]interface{})[0])

	if arrayInterfaceLen(o, []int{}) < 2 {
		return
	}

	el.LocationInInput = parseLocationInInput(o.([]interface{})[1])

	return
}

func parseLocationInGroonga(o interface{}) (lig LocationInGroonga) {
	lig.FunctionName = interfaceToString(o, 0)
	lig.SourceFileName = interfaceToString(o, 1)
	lig.LineNumber = interfaceToInt(o, 2)

	return
}

func parseLocationInInput(o interface{}) (lii LocationInInput) {
	lii.InputFileName = interfaceToString(o, 0)
	lii.LineNumber = interfaceToInt(o, 1)
	lii.LineContent = interfaceToString(o, 2)

	return
}
