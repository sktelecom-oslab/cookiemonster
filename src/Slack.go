package main

import (
	"github.com/bluele/slack"
)

const (
	webhookUrl  = "https://hooks.slack.com/services/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"
	apiToken    = "AIzaSyA-sd-Yb9aJepc0nh_LuWabpRRQoJelA3I"
	channelName = "jenkins"
)

func postSlack(message string) {
	webHook := slack.NewWebHook(webhookUrl)
	payload := &slack.WebHookPostPayload{Text: message, Channel: channelName}
	if err := webHook.PostMessage(payload); err != nil {
		panic(err)
	}
}
