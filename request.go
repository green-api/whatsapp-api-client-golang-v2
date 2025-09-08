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

type requestType struct {
	GetParams   string
	FormData    bool
	SetMimetype mtype
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

type mtype struct {
	Mimetype string
	FileName string
}

func WithSetMimetype(mtype mtype) requestOptions {
	return func(r *requestType) error {
		r.SetMimetype = mtype
		return nil
	}
}

func WithMediaHost(b bool) requestOptions {
	return func(r *requestType) error {
		r.MediaHost = b
		return nil
	}
}

func (a *GreenAPI) Request(HTTPMethod, APIMethod string, requestBody []byte, options ...requestOptions) (*APIResponse, error) {
	r := &requestType{}
	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	return a.request(HTTPMethod, APIMethod, r.GetParams, r.SetMimetype, r.FormData, r.MediaHost, requestBody)
}

func (a *GreenAPIPartner) PartnerRequest(HTTPMethod, APIMethod string, requestBody []byte) (*APIResponse, error) {
	client := &fasthttp.Client{}
	client.Name = "green-api-go-client " + a.Email

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("https://api.green-api.com/v3/partner/%s/%s", APIMethod, a.PartnerToken))

	req.Header.SetMethod(HTTPMethod)
	req.Header.Set("Content-Type", "application/json")

	if requestBody != nil {
		req.SetBody(requestBody)
	}

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

func MultipartRequest(method, url string, requestBody []byte) (*fasthttp.Request, error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	var filePath string

	var UnmarshaledBody map[string]interface{}

	err := json.Unmarshal(requestBody, &UnmarshaledBody)
	if err != nil {
		return nil, err
	}

	if v, ok := UnmarshaledBody["file"]; ok {
		filePath = v.(string)
	} else {
		return nil, fmt.Errorf("failed to retrieve FilePath from requestBody")
	}

	for key, value := range UnmarshaledBody {
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

	var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
		quoteEscaper.Replace("file"), quoteEscaper.Replace(filepath.Base(filePath))))

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

func (a *GreenAPI) request(HTTPMethod, APIMethod, GetParams string, SetMimetype mtype, FormData, MediaHost bool, requestBody []byte) (*APIResponse, error) {
	client := &fasthttp.Client{
		// Dial: func(addr string) (net.Conn, error) {
		//     return fasthttp.DialTimeout(addr, 10*time.Second)
		// },
		// ReadTimeout: time.Second * 10,
		// WriteTimeout: time.Second * 10,
	}
	client.Name = "green-api-go-client"

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%s/waInstance%s/%s/%s", a.APIURL, a.IDInstance, APIMethod, a.APITokenInstance))

	req.Header.SetMethod(HTTPMethod)
	req.Header.Set("Content-Type", "application/json")

	if MediaHost {
		req.SetRequestURI(fmt.Sprintf("%s/waInstance%s/%s/%s", a.MediaURL, a.IDInstance, APIMethod, a.APITokenInstance))
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

	if SetMimetype.Mimetype != "" {
		req.Header.SetContentType(SetMimetype.Mimetype)
		req.Header.Set("GA-Filename", SetMimetype.FileName)
	}

	if requestBody != nil {
		req.SetBody(requestBody)
	}

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
