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

type request struct {
	GetParams   bool
	FormData    bool
	SetMimetype string
	Partner     bool
	MediaHost   bool
}

type requestOptions func(*request)

// TODO: передавать сам addUrl вместо bool
func WithGetParams(b bool) requestOptions {
	return func(r *request) {
		r.GetParams = b
	}
}

func WithFormData(b bool) requestOptions {
	return func(r *request) {
		r.FormData = b
	}
}

func WithSetMimetype(mtype string) requestOptions {
	return func(r *request) {
		r.SetMimetype = mtype
	}
}

func WithPartner(b bool) requestOptions {
	return func(r *request) {
		r.Partner = b
	}
}

func WithMediaHost(b bool) requestOptions {
	return func(r *request) {
		r.MediaHost = b
	}
}

// TODO: добавить приватный request func
// посмотреть в python sdk
func (a *GreenAPI) Request(httpMethod, APImethod string, requestBody map[string]interface{}, options ...requestOptions) (any, error) {
	client := &fasthttp.Client{}

	r := &request{}
	for _, o := range options {
		o(r)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%s/waInstance%s/%s/%s", a.APIURL, a.IDInstance, APImethod, a.APITokenInstance))

	req.Header.SetMethod(httpMethod)

	if r.MediaHost {
		req.SetRequestURI(fmt.Sprintf("%s/waInstance%s/%s/%s", a.MediaURL, a.IDInstance, APImethod, a.APITokenInstance))
	}

	if r.Partner {
		req.SetRequestURI(fmt.Sprintf("%s/partner/%s/%s", a.APIURL, APImethod, a.PartnerToken))
	}

	if r.GetParams {
		var addUrl string
		if v, ok := requestBody["addUrl"]; ok {
			addUrl = v.(string)
		} else {
			return nil, fmt.Errorf("error while retreiving GET params and adding to URL")
		}
		req.SetRequestURI(req.URI().String() + addUrl)
	}

	if r.FormData {
		req, err := MultipartRequest(APImethod, req.URI().String(), requestBody)
		if err != nil {
			return nil, err
		}
		defer fasthttp.ReleaseRequest(req)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		if err := client.Do(req, resp); err != nil {
			return nil, fmt.Errorf("request error: %s", err)
		}
		fmt.Println(req.URI())
		return &APIResponse{
			StatusCode: resp.StatusCode(),
			Body:       resp.Body(),
			Timestamp:  time.Now().Format("15:04:05.000"),
		}, nil
	}

	if r.SetMimetype != "" {
		req.Header.SetContentType(r.SetMimetype)
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error when serializing data to JSON: %s", err)
	}
	req.SetBody([]byte(jsonData))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	fmt.Println(req.URI())
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("request error: %s", err)
	}

	return &APIResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Timestamp:  time.Now().Format("15:04:05.000"),
	}, nil
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

	for key, value := range requestBody {
		if key == "file" {
			continue
		}
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

	req.Header.SetMethod("POST")

	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.SetBody(buffer.Bytes())

	return req, nil
}

// func (a GreenAPI) NetHttpRequest(httpMethod, APImethod string, requestBody map[string]interface{}) (any, error) {
// 	client := &http.Client{}

// 	var reqBody io.Reader
// 	if requestBody != nil {
// 		jsonData, err := json.Marshal(requestBody)
// 		if err != nil {
// 			return nil, fmt.Errorf("ошибка при сериализации данных в JSON: %s", err)
// 		}
// 		reqBody = bytes.NewBuffer(jsonData)
// 	}

// 	req, err := http.NewRequest(httpMethod, a.getRequestURL(APImethod), reqBody)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при создании запроса: %s", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при запросе: %s", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при чтении тела ответа: %s", err)
// 	}

// 	return &ApiResponse{
// 		StatusCode: resp.StatusCode,
// 		Body:       string(body),
// 		Timestamp:  time.Now().Format("15:04:05.000"),
// 	}, nil
// }
