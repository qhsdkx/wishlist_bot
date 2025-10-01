package bot

import (
	"strings"
)

type CallbackData struct {
	mode string
	id   string
	page string
}

func NewCallbackData(mode, id, page string) CallbackData {
	return CallbackData{
		mode: mode,
		id:   id,
		page: page,
	}
}

func (c CallbackData) string() string {
	return c.mode + "|" + c.id + "|" + c.page
}

func (c CallbackData) Mode() string {
	return c.mode
}

func (c CallbackData) Id() string {
	return c.id
}

func (c CallbackData) Page() string {
	return c.page
}

func parseCallback(value string) CallbackData {
	data := strings.Split(value, "|")
	return CallbackData{
		mode: data[0],
		id:   data[1],
		page: data[2],
	}
}
