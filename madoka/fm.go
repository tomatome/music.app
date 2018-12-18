// fm.go
package madoka

// fm 信息
func FmInfo(id string) (string, error) {
	preParams := "{\"id\": \"" + id + "\", \"csrf_token\": \"\"}"
	params, encSecKey, err := EncParams(preParams)
	if err != nil {
		return "", err
	}
	res, resErr := post("http://music.163.com/weapi/djradio/get", params, encSecKey)
	if resErr != nil {
		return "", resErr
	}
	return res, nil
}

// fmCatalogue - get fm list
func FmCatalogue() (string, error) {
	preParams := "{\"csrf_token\": \"\"}"
	params, encSecKey, encErr := EncParams(preParams)
	if encErr != nil {
		return "", encErr
	}
	res, resErr := post("http://music.163.com/weapi/djradio/category/get", params, encSecKey)
	if resErr != nil {
		return "", resErr
	}
	return res, nil
}

func FmHotList(category, id string, page, limit int) (string, error) {
	_offset, _limit := formatParams(page, limit)
	preParams := "{\"cat\":\"" + "全部" +
		"\", \"cateId\":\"" + id +
		"\", \"type\":\"" + id +
		"\", \"categoryId\":\"" + id +
		"\", \"category\":\"" + category +
		"\", \"offset\":\"" + _offset +
		"\", \"limit\":\"" + _limit +
		"\"}"

	params, encSecKey, err := EncParams(preParams)
	if err != nil {
		return "", err
	}
	res, resErr := post("http://music.163.com/weapi/djradio/hot/v1", params, encSecKey)
	if resErr != nil {
		return "", resErr
	}
	return res, nil
}
