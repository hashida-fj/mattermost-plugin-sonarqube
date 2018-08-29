package main

import (
	// "bytes"
	// "strings"
    //  "text/template"
	// "errors"

	"github.com/mattermost/mattermost-server/model"
)

type WebhookResponse struct {
	AnalysedAt string `json:"analysedAt"`
	Project    struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"project"`
	Properties struct {
	} `json:"properties"`
	QualityGate struct {
		Conditions []struct {
			ErrorThreshold string `json:"errorThreshold"`
			Metric         string `json:"metric"`
			OnLeakPeriod   bool   `json:"onLeakPeriod"`
			Operator       string `json:"operator"`
			Status         string `json:"status"`
			Value          string `json:"value,omitempty"`
		} `json:"conditions"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"qualityGate"`
	ServerURL string `json:"serverUrl"`
	Status    string `json:"status"`
	TaskID    string `json:"taskId"`
}


func (w *WebhookResponse) SlackAttachment() (*model.SlackAttachment, error) {

	return &model.SlackAttachment{
		Color: "#95b7d0",

		Text: w.QualityGate.Status + w.QualityGate.Conditions[3].Value + ":" + w.QualityGate.Conditions[3].Status,
		Title: "sonar",

		AuthorName: "SonarQube",
		AuthorLink: w.ServerURL,
	}, nil

}
