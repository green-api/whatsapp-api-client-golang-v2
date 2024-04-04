package greenapi

type AccountCategory struct {
	GreenAPI GreenAPIInterface
}

func (c AccountCategory) GetSettings() (interface{}, error) {
	return c.GreenAPI.Request("GET", "getSettings", nil)
}

func (c AccountCategory) GetInstances() (interface{}, error) {
	return c.GreenAPI.Request("GET", "getInstances", map[string]interface{}{}) //обязательно передавать пустой интерфейс вместо nil
}

func (c AccountCategory) SetSettings(requestBody map[string]interface{}) (interface{}, error) {
	return c.GreenAPI.Request("POST", "setSettings", requestBody)
}
