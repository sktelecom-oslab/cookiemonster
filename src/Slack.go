/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
