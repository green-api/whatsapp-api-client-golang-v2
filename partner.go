package greenapi

import "encoding/json"

type PartnerCategory struct {
	GreenAPIPartner GreenAPIPartnerInterface
}

// ------------------------------------------------------------------ GetInstances block

func (c PartnerCategory) GetInstances() (*APIResponse, error) {
	return c.GreenAPIPartner.PartnerRequest("GET", "getInstances", nil)
}

// ------------------------------------------------------------------ CreateInstance block

func (c PartnerCategory) CreateInstance(options ...SetSettingsOption) (*APIResponse, error) {
	r := &RequestSetSettings{}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPIPartner.PartnerRequest("POST", "createInstance", jsonData)
}

// ------------------------------------------------------------------ DeleteInstanceAccount block

type RequestDeleteInstanceAccount struct {
	IdInstance int `json:"idInstance"`
}

func (c PartnerCategory) DeleteInstanceAccount(idInstance int) (*APIResponse, error) {
	r := &RequestDeleteInstanceAccount{
		IdInstance: idInstance,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPIPartner.PartnerRequest("POST", "deleteInstanceAccount", jsonData)
}
