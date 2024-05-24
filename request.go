package greenapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/valyala/fasthttp"
)

// func (a GreenAPI) Webhook() GreenAPIWebhook {
// 	return GreenAPIWebhook{
// 		GreenAPI: a,

// 		ErrorChannel: make(chan error),
// 	}
// }

type requestType struct {
	GetParams   string
	FormData    bool
	SetMimetype string
	Partner     bool
	MediaHost   bool
}

type requestOptions func(*requestType) error

func WithGetParams(addUrl string) requestOptions {
	return func(r *requestType) error {
		r.GetParams = addUrl
		return nil
	}
}

func WithFormData(b bool) requestOptions {
	return func(r *requestType) error {
		r.FormData = b
		return nil
	}
}

func WithSetMimetype(mtype string) requestOptions {
	return func(r *requestType) error {
		r.SetMimetype = mtype
		return nil
	}
}

func WithPartner(b bool) requestOptions {
	return func(r *requestType) error {
		r.Partner = b
		return nil
	}
}

func WithMediaHost(b bool) requestOptions {
	return func(r *requestType) error {
		r.MediaHost = b
		return nil
	}
}

func (a *GreenAPI) Request(HTTPMethod, APIMethod string, requestBody map[string]interface{}, options ...requestOptions) (*APIResponse, error) {
	r := &requestType{}
	for _, o := range options {
		o(r)
	}

	return a.request(HTTPMethod, APIMethod, r.GetParams, r.SetMimetype, r.FormData, r.Partner, r.MediaHost, requestBody)
}

func MultipartRequest(method, url string, requestBody map[string]interface{}) (*fasthttp.Request, error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	var filePath string

	if v, ok := requestBody["file"]; ok {
		filePath = v.(string)
	} else {
		return nil, fmt.Errorf("failed to retrieve FilePath from requestBody")
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
	defer file.Close()

	//this is modified code of writer.CreateFormFile function
	//the original function does not allow to set Content-Type of a particular field of Form-Data other than application/octet
	h := make(textproto.MIMEHeader)

	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
		escapeQuotes("file"), escapeQuotes(filepath.Base(filePath))))

	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		return nil, err
	}
	h.Set("Content-Type", mtype.String())

	part, err := writer.CreatePart(h)
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

// TODO: добавить нормальную обработку ошибок
func (a *GreenAPI) request(HTTPMethod, APIMethod, GetParams, SetMimetype string, FormData, Partner, MediaHost bool, requestBody map[string]interface{}) (*APIResponse, error) {
	client := &fasthttp.Client{}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%s/waInstance%s/%s/%s", a.APIURL, a.IDInstance, APIMethod, a.APITokenInstance))

	req.Header.SetMethod(HTTPMethod)

	if MediaHost {
		req.SetRequestURI(fmt.Sprintf("%s/waInstance%s/%s/%s", a.MediaURL, a.IDInstance, APIMethod, a.APITokenInstance))
	}

	if Partner {
		req.SetRequestURI(fmt.Sprintf("%s/partner/%s/%s", a.APIURL, APIMethod, a.PartnerToken))
	}

	if GetParams != "" {
		req.SetRequestURI(req.URI().String() + GetParams)
	}

	if FormData {
		req, err := MultipartRequest(APIMethod, req.URI().String(), requestBody)
		if err != nil {
			return nil, err
		}
		defer fasthttp.ReleaseRequest(req)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		if err := client.Do(req, resp); err != nil {
			return nil, fmt.Errorf("request error: %s", err)
		}

		return &APIResponse{
			StatusCode:    resp.StatusCode(),
			StatusMessage: resp.Header.StatusMessage(),
			Body:          resp.Body(),
			Timestamp:     time.Now(),
		}, nil
	}

	if SetMimetype != "" {
		req.Header.SetContentType(SetMimetype)
	}

	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("error when serializing data to JSON: %s", err)
		}
		req.SetBody([]byte(jsonData))
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// fmt.Println(req.URI())
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("request error: %s", err)
	}

	//fmt.Println(req.Body())

	return &APIResponse{
		StatusCode:    resp.StatusCode(),
		StatusMessage: resp.Header.StatusMessage(),
		Body:          resp.Body(),
		Timestamp:     time.Now(),
	}, nil
}

// TODO: figure out what to do with this
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
