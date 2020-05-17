package main

import (
	"encoding/json"
	"html/template"
)

type News struct {
	Headline    string
	Description string
	Text        template.HTML
}

func ConvertToNews(msg json.RawMessage) (interface{}, error) {
	var news News
	err := json.Unmarshal([]byte(msg), &news)
	if err != nil {
		return nil, err
	}
	return news, nil
}
