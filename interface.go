package greenapi

type GreenAPIInterface interface {
	Request(httpMethod, APImethod string, data map[string]interface{}) (any, error)
	//OldRequest(httpMethod, APImethod string, requestBody map[string]interface{}) (any, error)
}
