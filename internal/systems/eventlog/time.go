package eventlog

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lestrrat-go/strftime"
)

// AddTimeToEmbed formats the current time using timeFmt and adds it to e.
// The timeFmt can either be "discord", "juche", or a strftime string.
func AddTimeToEmbed(timeFmt string, e *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	t := time.Now().In(time.UTC)
	switch timeFmt {
	case "discord":
		e.Timestamp = t.Format(time.RFC3339)
	case "juche":
		e.Footer = &discordgo.MessageEmbedFooter{Text: formatJuche(t)}
	default:
		e.Footer = &discordgo.MessageEmbedFooter{Text: format(timeFmt, t)}
	}
	return e
}

// formatJuche formats the given time in Juche calendar format
func formatJuche(t time.Time) string {
	return fmt.Sprintf("%02d:%02d %02d-%02d Juche %d", t.Hour(), t.Minute(), t.Day(), t.Month(), t.Year()-1911)
}

// format formats t using timeFmt
func format(timeFmt string, t time.Time) string {
	timeStr, err := strftime.Format(timeFmt, t)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return timeStr
}
