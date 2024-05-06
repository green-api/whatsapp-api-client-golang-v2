package greenapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

// func (a GreenAPI) Webhook() GreenAPIWebhook {
// 	return GreenAPIWebhook{
// 		GreenAPI: a,

// 		ErrorChannel: make(chan error),
// 	}
// }

func (a GreenAPI) Request(httpMethod, APImethod string, requestBody map[string]interface{}) (any, error) {

	client := &fasthttp.Client{}

	if APImethod == "sendFileByUpload" {
		req, err := MultipartRequest(APImethod, a.getRequestURL(APImethod), requestBody)
		if err != nil {
			return nil, err
		}
		defer fasthttp.ReleaseRequest(req)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		if err := client.Do(req, resp); err != nil {
			return nil, fmt.Errorf("ошибка при запросе: %s", err)
		}

		// var response interface{}
		// err = json.Unmarshal(resp.Body(), &response)
		// if err != nil {
		// 	color.Green("Body: %s", req.Body())
		// 	return nil, fmt.Errorf("error while unmarshal byte response: %s", err)
		// }

		return &ApiResponse{
			StatusCode: resp.StatusCode(),
			Body:       string(resp.Body()),
			Timestamp:  time.Now().Format("15:04:05.000"),
		}, nil
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(a.getRequestURL(APImethod))

	req.Header.SetMethod(httpMethod)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации данных в JSON: %s", err)
	}
	req.SetBody([]byte(jsonData))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		fmt.Println(req.URI())
		return nil, fmt.Errorf("ошибка при запросе: %s", err)
	}

	// var response interface{}
	// err = json.Unmarshal(resp.Body(), &response)
	// if err != nil {
	// 	color.Green("Body: %s", req.Body())
	// 	return nil, fmt.Errorf("error while unmarshal byte response: %s", err)
	// }

	return &ApiResponse{
		StatusCode: resp.StatusCode(),
		Body:       string(resp.Body()),
		Timestamp:  time.Now().Format("15:04:05.000"),
	}, nil
}

func (a GreenAPI) getRequestURL(APIMethod string) string {
	switch APIMethod {
	case "createInstance", "deleteInstanceAccount", "getInstances":
		return fmt.Sprintf("%s/partner/%s/%s", a.Host, APIMethod, a.PartnerToken)
	case "sendFileByUpload":
		return fmt.Sprintf("%s/waInstance%s/%s/%s", a.MediaHost, a.IDInstance, APIMethod, a.APITokenInstance)
	default:
		return fmt.Sprintf("%s/waInstance%s/%s/%s", a.Host, a.IDInstance, APIMethod, a.APITokenInstance)
	}
}

func MultipartRequest(method, url string, requestBody map[string]interface{}) (*fasthttp.Request, error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	var filePath string

	if v, ok := requestBody["file"]; ok {
		filePath = v.(string)
	} else {
		return nil, fmt.Errorf("failed to retreive FilePath from requestBody")
	}

	delete(requestBody, "file")

	for key, value := range requestBody {
		err := writer.WriteField(key, value.(string))
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()

	req.SetRequestURI(url)
	fmt.Println(req.URI())

	req.Header.SetMethod("POST")

	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.SetBody(buffer.Bytes())

	return req, nil
}