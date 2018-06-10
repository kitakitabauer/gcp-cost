package main

import (
	nSlack "github.com/nlopes/slack"
	"sort"
	"strconv"
)

type Type struct {
	Cost       float64
	Attachment nSlack.Attachment
}

const (
	slackToken     = "" // your token
	slackChannel   = "" // your channel
	slackIconEmoji = "" // your icon emoji

	thresoldOfDisplayUSD = 0.01
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

func init() {
	slackCli = nSlack.New(slackToken)
}

func divideType(costs []CostResult) map[string][]CostResult {
	m := make(map[string][]CostResult, len(costs))
	for _, cost := range costs {
		if _, found := m[cost.ResourceType]; !found {
			m[cost.ResourceType] = []CostResult{cost}
		} else {
			m[cost.ResourceType] = append(m[cost.ResourceType], cost)
		}
	}
	return m
}

func float64ToString(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func calcTypeCost(service string, costs []CostResult) Type {
	var fields []nSlack.AttachmentField
	var typeCost float64
	for _, cost := range costs {
		if cost.Cost <= thresoldOfDisplayUSD {
			// do not display because it's too low price
			sumCost += cost.Cost
			continue
		}

		field := nSlack.AttachmentField{
			Title: cost.ResourceType,
			Value: float64ToString(cost.Cost, 3) + " (" + cost.Currency + ")",
			Short: true,
		}
		fields = append(fields, field)

		sumCost += cost.Cost
		typeCost += cost.Cost
	}

	return Type{
		Cost: typeCost,
		Attachment: nSlack.Attachment{
			Title:  "[" + service + "]",
			Fields: fields,
		},
	}
}

func sortCost(groups []Type) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Cost > groups[j].Cost
	})
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

func sendToSlack(targetYmd string, costs []CostResult) error {
	types := divideType(costs)
	typeCosts := make([]Type, 0, len(types))
	for name, t := range types {
		_type := calcTypeCost(name, t)
		if _type.Cost == 0 {
			continue
		}
		typeCosts = append(typeCosts, _type)
	}

	sortCost(typeCosts)

	var colorIndex int
	attachments := make([]nSlack.Attachment, 0, len(typeCosts))
	for _, _type := range typeCosts {
		_type.Attachment.Color = selectColor(colorIndex)
		colorIndex = selectColorIndex(colorIndex)

		attachments = append(attachments, _type.Attachment)
	}

	msg := nSlack.PostMessageParameters{
		Username:    "GCP Cost of " + targetYmd,
		Channel:     slackChannel,
		IconEmoji:   slackIconEmoji,
		Attachments: attachments,
	}

	total := "-------------------------------------\nTotal cost: " +
		float64ToString(sumCost, 6) +
		" (" + costs[0].Currency +
		")\n-------------------------------------"
	_, _, err := slackCli.PostMessage(msg.Channel, total, msg)
	return err
}
