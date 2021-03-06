package groupmestatsbot

import (
	"fmt"

	"github.com/MagnusFrater/groupme"
)

// SprintTotalMessages formats an Total Messages Bot post and returns the resulting string.
func (s *Stats) SprintTotalMessages() string {
	return fmt.Sprintf("Total Messages: %d messages", s.TotalMessages())
}

// SprintAverageMessageLength formats an Average Message Length Bot post and returns the resulting string.
func (s *Stats) SprintAverageMessageLength() string {
	averageMessageLength := s.AverageMessageLength()
	if averageMessageLength == -1 {
		averageMessageLength = 0
	}

	return fmt.Sprintf("Average Message Length: %d words", averageMessageLength)
}

// SprintTopMessages formats a Top Messages Bot post and returns the resulting string.
func (s *Stats) SprintTopMessages(limit int) string {
	str := fmt.Sprintf("Top Messages\n%s\n", messageDivider)

	topMessages := s.TopMessages(limit)
	if len(topMessages) == 0 {
		str += "\nThere are no messages."
		return str
	}

	for i, message := range topMessages {
		str += fmt.Sprintf("%d) (%d) %s:\n", i+1, len(message.FavoritedBy), message.Name)

		if len(message.Attachments) > 0 {
			for i, attachment := range message.Attachments {
				switch attachment.Type {
				case groupme.ImageAttachment:
					str += fmt.Sprintf("image: %s", attachment.URL)
				}

				// only put newline if there are more attachments, or if the message has text
				if i < len(message.Attachments)-1 || message.Text != "" {
					str += "\n"
				}
			}
		}

		if message.Text != "" {
			str += fmt.Sprintf("\"%s\"", message.Text)
		}

		// don't put newline after last ranking
		if i < len(topMessages)-1 {
			str += "\n\n"
		}
	}

	return str
}
