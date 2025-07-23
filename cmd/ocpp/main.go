package main

import "github.com/squishmeist/ocpp-go/service/ocpp"

const (
	topicName        = "topic.1"
	subscriptionName = "subscription.1"
)


func main() {
state := &ocpp.State{}
err := ocpp.ListenToTopicAndProcess(state, topicName, subscriptionName)
if err != nil {
    panic(err)
}
}
