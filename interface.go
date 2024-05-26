package greenapi

type GreenAPIInterface interface {
	Request(httpMethod, APImethod string, requestBody []byte, options ...requestOptions) (*APIResponse, error)
	//NetHttpRequest(httpMethod, APImethod string, requestBody map[string]interface{}) (any, error)
}
