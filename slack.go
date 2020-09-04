package main

import (
	"strconv"

	"github.com/nlopes/slack"
	nSlack "github.com/nlopes/slack"
)

var (
	sumCost  float64
	slackCli *nSlack.Client

	googleColors = []string{
		"#4285F4", // blue
		"#0F9D58", // green
		"#F4B400", // yellow
		"#DB4437", // red
	}
)

const slackToken = "" // your token

func init() {
	slackCli = nSlack.New(slackToken)
}

func float64ToString(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func selectColor(idx int) string {
	return googleColors[idx]
}

func selectColorIndex(idx int) int {
	idx++
	if idx == len(googleColors) {
		return 0
	}
	return idx
}

func sendToSlack(targetYmd, channel string, costs []CostResult) error {
	var colorIndex int
	attachments := make([]nSlack.Attachment, 0, len(costs)+1)
	for i := range costs {
		if costs[i].Cost == 0 {
			continue
		}

		attachment := nSlack.Attachment{
			Color: selectColor(colorIndex),
			Fields: []nSlack.AttachmentField{
				{
					Title: costs[i].Service,
					Value: float64ToString(costs[i].Cost, 3) + " (USD)",
					Short: true,
				},
			},
		}
		attachments = append(attachments, attachment)
		sumCost += costs[i].Cost

		colorIndex = selectColorIndex(colorIndex)
	}

	total := "-------------------------------------\nTotal cost: " +
		float64ToString(sumCost, 3) +
		" (USD)\n-------------------------------------"

	attachments = append([]nSlack.Attachment{{Text: total}}, attachments[0:]...)

	msg := nSlack.PostMessageParameters{
		Username:  "GCP Cost of " + targetYmd,
		Channel:   channel,
		IconEmoji: ":gcp:",
	}

	_, _, err := slackCli.PostMessage(
		msg.Channel,
		slack.MsgOptionPostMessageParameters(msg),
		slack.MsgOptionAttachments(attachments...),
	)
	return err
}
