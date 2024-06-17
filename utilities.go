package greenapi

import (
	"fmt"
	"net/url"
	"strings"
)

func ValidateChatId(chatId string) error {
	if !strings.HasSuffix(chatId, "@c.us") && !strings.HasSuffix(chatId, "@g.us") {
		return fmt.Errorf("chatId must end with \"@c.us\" or \"@g.us\"")
	}
	return nil
}

func ValidateMessageLength(message string, limit int) error {
	if len(message) > limit {
		return fmt.Errorf("length of the message exceeds the limit of %v", limit)
	}
	return nil
}

func ValidateURL(link string) error {
	_, err := url.Parse(link)
	if err!= nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}
	return nil
}