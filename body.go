package gogroongahttpdmanage

func parseBodySelect(body interface{}) (bs BodySelect) {
	length := len(body.([]interface{}))

	bs.SearchResults = make([]SearchResult, 0, length)

	for i := 0; i < length; i++ {
		o := body.([]interface{})[i]
		var sr SearchResult
		sr.NHits = arrayInterfaceToInt(o, []int{0, 0})

		// カラム情報
		sr.Columns = make([]GroongaType, 0, 8)
		for i := 0; i < len(o.([]interface{})[1].([]interface{})); i++ {
			item := o.([]interface{})[1].([]interface{})[i]

			name := item.([]interface{})[0].(string)
			typ := item.([]interface{})[1].(string)
			gt := GroongaType{
				Name: name,
				Type: Type(typ),
			}
			sr.Columns = append(sr.Columns, gt)
		}

		// レコード
		recordLen := len(o.([]interface{})[2].([]interface{}))
		sr.Records = make([]interface{}, 0, recordLen)
		for i := 0; i < recordLen; i++ {
			index := i + 2

			obj := make([]interface{}, 0, len(sr.Columns))
			for k := 0; k < len(sr.Columns); k++ {
				obj = append(obj, o.([]interface{})[index].([]interface{})[k])
			}
			sr.Records = append(sr.Records, obj)
		}
		bs.SearchResults = append(bs.SearchResults, sr)
	}

	return bs
}
