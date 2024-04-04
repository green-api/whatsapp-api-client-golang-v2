package greenapi

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

// func (a GreenAPI) Webhook() GreenAPIWebhook {
// 	return GreenAPIWebhook{
// 		GreenAPI: a,

// 		ErrorChannel: make(chan error),
// 	}
// }

func (a GreenAPI) Request(httpMethod, APImethod string, requestBody map[string]interface{}) (interface{}, error) {
	client := &fasthttp.Client{}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	if APImethod == "createInstance" || APImethod == "deleteInstanceAccount" || APImethod == "getInstances" {
		req.SetRequestURI(fmt.Sprintf("%s/partner/%s/%s", a.Host, APImethod, a.PartnerToken))
	} else {
		req.SetRequestURI(fmt.Sprintf("%s/waInstance%s/%s/%s", a.Host, a.IDInstance, APImethod, a.APITokenInstance))
	}

	req.Header.SetMethod(httpMethod)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации данных в JSON: %s", err)
	}
	req.SetBody([]byte(jsonData))

	// TODO: handle fileUpload cases. Handle filePath inside requestBody

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("ошибка при запросе: %s", err)
	}

	var response interface{}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshal byte response: %s", err)
	}

	return response, nil
}
