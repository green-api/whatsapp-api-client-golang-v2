package greenapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testLID = "123456789012345@lid"

func TestSendToLID(t *testing.T) {
	tests := []struct {
		endpoint string
		send     func(SendingCategory) (*APIResponse, error)
	}{
		{"sendMessage", func(sending SendingCategory) (*APIResponse, error) {
			return sending.SendMessage(testLID, "Choose a language")
		}},
		{"sendFileByUrl", func(sending SendingCategory) (*APIResponse, error) {
			return sending.SendFileByUrl(testLID, "https://example.test/menu.png", "menu.png")
		}},
	}
	for _, tt := range tests {
		t.Run(tt.endpoint, func(t *testing.T) {
			wantPath := "/waInstance123/" + tt.endpoint + "/test-token"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.RequestURI != wantPath {
					t.Errorf("got %s %s, want POST %s", r.Method, r.RequestURI, wantPath)
				}
				var payload struct {
					ChatID string `json:"chatId"`
				}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Errorf("decode request: %v", err)
				}
				if payload.ChatID != testLID {
					t.Errorf("got chatId %q, want %q", payload.ChatID, testLID)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()
			api := GreenAPI{APIURL: server.URL, IDInstance: "123", APITokenInstance: "test-token"}
			response, err := tt.send(api.Sending())
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != http.StatusNoContent {
				t.Errorf("got HTTP %d, want %d", response.StatusCode, http.StatusNoContent)
			}
		})
	}
}

func TestValidateChatId(t *testing.T) {
	if err := ValidateChatId("123@c.us", "123@g.us", testLID); err != nil {
		t.Errorf("valid chat IDs: %v", err)
	}
	for _, invalid := range []string{"", "123", "123@unknown", testLID + ".extra"} {
		if err := ValidateChatId("123@c.us", invalid); err == nil {
			t.Errorf("accepted invalid chat ID %q", invalid)
		}
	}
}

type chatIDRequestFunc func(method, endpoint string, body []byte) (*APIResponse, error)

func (request chatIDRequestFunc) Request(method, endpoint string, body []byte, _ ...requestOptions) (*APIResponse, error) {
	return request(method, endpoint, body)
}

func TestAddGroupParticipantByLID(t *testing.T) {
	const groupID = "123@g.us"
	called := false
	groups := GroupsCategory{GreenAPI: chatIDRequestFunc(func(method, endpoint string, body []byte) (*APIResponse, error) {
		called = true
		if method != http.MethodPost || endpoint != "addGroupParticipant" {
			t.Errorf("got %s %s, want POST addGroupParticipant", method, endpoint)
		}
		var payload RequestModifyGroupParticipant
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload.GroupId != groupID || payload.ParticipantChatId != testLID {
			t.Errorf("got groupId %q and participantChatId %q, want %q and %q", payload.GroupId, payload.ParticipantChatId, groupID, testLID)
		}
		return &APIResponse{StatusCode: http.StatusOK}, nil
	})}
	if _, err := groups.AddGroupParticipant(groupID, testLID); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("request was not sent")
	}
}

func TestChatIDValidationBeforeRequest(t *testing.T) {
	request := chatIDRequestFunc(func(method, endpoint string, _ []byte) (*APIResponse, error) {
		t.Fatalf("unexpected request: %s %s", method, endpoint)
		return nil, nil
	})
	groups := GroupsCategory{GreenAPI: request}
	contacts := ContactsCategory{GreenAPI: request}
	sending := SendingCategory{GreenAPI: request}
	statuses := StatusesCategory{GreenAPI: request}
	if _, err := groups.GetGroupData(testLID); err == nil {
		t.Error("GetGroupData accepted a LID as groupId")
	}
	if _, err := groups.AddGroupParticipant("123@c.us", testLID); err == nil {
		t.Error("AddGroupParticipant accepted a phone chat ID as groupId")
	}
	if _, err := groups.CreateGroup("Group", []string{"123@g.us"}); err == nil {
		t.Error("CreateGroup accepted a group as a participant")
	}
	if _, err := contacts.AddContact(testLID, "Name"); err == nil {
		t.Error("AddContact accepted a LID as a phone chat ID")
	}
	if _, err := sending.SendMessage("123@unknown", "Hello"); err == nil {
		t.Error("SendMessage accepted an unknown suffix")
	}
	if _, err := statuses.SendTextStatus("Hello", OptionalParticipantsTextStatus([]string{"123@g.us"})); err == nil {
		t.Error("SendTextStatus accepted a group as a recipient")
	}
}
