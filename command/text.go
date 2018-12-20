package command

import "encoding/json"

type Text struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

func JsonEncode(argv ...interface{}) (string, bool) {
	var texts []Text
	for _, v := range argv {
		switch t := v.(type) {
		case []Text:
			return _jsonEncode(t...), true
		case Text:
			texts = append(texts, t)
		default:
			return "", false
		}
	}
	return _jsonEncode(texts...), true
}

func _jsonEncode(argv ...Text) string {
	var texts []Text
	for _, v := range argv {
		texts = append(texts, v)
	}
	res, _ := json.Marshal(texts)
	return string(res)
}
