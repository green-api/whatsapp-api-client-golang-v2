package greenapi

type GreenAPIInterface interface {
	Request(httpMethod, APImethod string, data map[string]interface{}, options ...requestOptions) (*APIResponse, error)
	//NetHttpRequest(httpMethod, APImethod string, requestBody map[string]interface{}) (any, error)
}
