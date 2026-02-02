package bot

import (
	"strings"
)

const (
	action = iota
	mode
	id
	page
)

type CallbackData struct {
	action string
	mode   string
	id     string
	page   string
}

func NewCallbackData(action, mode, id, page string) CallbackData {
	return CallbackData{
		action: action,
		mode:   mode,
		id:     id,
		page:   page,
	}
}

func EmptyCallbackData() CallbackData {
	return CallbackData{}
}

func (c CallbackData) string() string {
	return c.action + "|" + c.mode + "|" + c.id + "|" + c.page
}

func (c CallbackData) Action() string {
	return c.action
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

func (c *CallbackData) SetAction(action string) {
	c.action = action
}

func (c *CallbackData) SetMode(mode string) {
	c.mode = mode
}

func (c *CallbackData) SetId(id string) {
	c.id = id
}

func (c *CallbackData) SetPage(page string) {
	c.page = page
}

func parseCallback(value string) CallbackData {
	data := strings.Split(value, "|")
	return CallbackData{
		action: data[action],
		mode:   data[mode],
		id:     data[id],
		page:   data[page],
	}
}
