package main

import (
	"github.com/bluele/slack"
)

const (
	webhookUrl = "https://hooks.slack.com/services/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"
	//token       = "xoxp-30956645507-135417633216-187836090932-4aac59518a595038412dcefa0468f380"
	channelName = "@wil"
)

func postSlack(message string) {
	webHook := slack.NewWebHook(webhookUrl)
	payload := &slack.WebHookPostPayload{Text: message}
	if err := webHook.PostMessage(payload); err != nil {
		panic(err)
	}
}
