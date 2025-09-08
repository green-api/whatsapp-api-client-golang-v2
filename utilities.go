package greenapi

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ValidateChatId(chatId ...string) error {

	for _, chat := range chatId {
		// Check for lid
		_, err := strconv.Atoi(chat)
		if err != nil {
			// Not lid, using old check
			if !strings.HasSuffix(chat, "@c.us") && !strings.HasSuffix(chat, "@g.us") {
				return fmt.Errorf("chatId must end with \"@c.us\" or \"@g.us\"\ngot %s instead", chat)
			}
		}
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
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}
	return nil
}
