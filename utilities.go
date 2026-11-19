package greenapi

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

// ValidateChatId checks the suffix of each WhatsApp chat ID.
func ValidateChatId(chatId ...string) error {
	return validateChatIDSuffixes(chatId, "@c.us", "@g.us", "@lid")
}

func validateGroupID(groupID string) error {
	return validateChatIDSuffixes([]string{groupID}, "@g.us")
}

func validatePersonalChatIDs(chatIDs ...string) error {
	return validateChatIDSuffixes(chatIDs, "@c.us", "@lid")
}

func validatePhoneChatID(chatID string) error {
	return validateChatIDSuffixes([]string{chatID}, "@c.us")
}

func validateChatIDSuffixes(chatIDs []string, suffixes ...string) error {
	for _, chatID := range chatIDs {
		at := strings.LastIndexByte(chatID, '@')
		if at < 0 || !slices.Contains(suffixes, chatID[at:]) {
			return fmt.Errorf("chatId must end with %s\ngot %s instead", strings.Join(suffixes, " or "), chatID)
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
	if err!= nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}
	return nil
}