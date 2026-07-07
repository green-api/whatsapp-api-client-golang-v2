package greenapi

import (
	"encoding/json"
)

type ContactsCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ AddContact

type RequestAddContact struct {
	ChatId            string `json:"chatId"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName,omitempty"`
	SaveInAddressbook bool   `json:"saveInAddressbook,omitempty"`
}

type AddContactOption func(*RequestAddContact) error

// Contact's last name.
func OptionalAddLastName(lastName string) AddContactOption {
	return func(r *RequestAddContact) error {
		r.LastName = lastName
		return nil
	}
}

// SaveInAddressbook. If present, the contact will be saved in the address book.
func OptionalAddSaveInAddressbook(saveInAddressbook bool) AddContactOption {
	return func(r *RequestAddContact) error {
		r.SaveInAddressbook = saveInAddressbook
		return nil
	}
}

func (c ContactsCategory) AddContact(chatId, firstName string, options ...AddContactOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestAddContact{
		ChatId:    chatId,
		FirstName: firstName,
	}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "addContact", jsonData)
}

// ------------------------------------------------------------------ EditContact

type RequestEditContact struct {
	ChatId            string `json:"chatId"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName,omitempty"`
	SaveInAddressbook bool   `json:"saveInAddressbook,omitempty"`
}

type EditContactOption func(*RequestEditContact) error

// Contact's last name.
func OptionalEditLastName(lastName string) EditContactOption {
	return func(r *RequestEditContact) error {
		r.LastName = lastName
		return nil
	}
}

// SaveInAddressbook. If present, the contact will be saved in the address book.
func OptionalEditSaveInAddressbook(saveInAddressbook bool) EditContactOption {
	return func(r *RequestEditContact) error {
		r.SaveInAddressbook = saveInAddressbook
		return nil
	}
}

func (c ContactsCategory) EditContact(chatId, firstName string, options ...EditContactOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestEditContact{
		ChatId:    chatId,
		FirstName: firstName,
	}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "editContact", jsonData)
}

// ------------------------------------------------------------------ DeleteContact

type RequestDeleteContact struct {
	ChatId string `json:"chatId"`
}

func (c ContactsCategory) DeleteContact(chatId string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestDeleteContact{
		ChatId: chatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "deleteContact", jsonData)
}
