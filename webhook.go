package main

import (
	// "strings"
    "bytes"
    "text/template"
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


func (w *WebhookResponse) buildTable() (string, error) {
    var rows bytes.Buffer
    //  var res = ""

    header := "|品質指標|現在の値|閾値|\n|:-|:-|:-|\n"
    tmplstr := "|{{.Status}}{{.Metric}}|{{.Value}}|{{.ErrorThreshold}}|\n"
    rowtmpl, err:= template.New("row").Parse(tmplstr)
    if err != nil { return "", err}

    rows.WriteString(header)
    for _, cond := range w.QualityGate.Conditions {
        if err := rowtmpl.Execute(&rows, cond); err != nil {
            return "error", err
        }
    }

    return rows.String(), nil
}

func (w *WebhookResponse) SlackAttachment() (*model.SlackAttachment, error) {

    table, err := w.buildTable()
    if err != nil {
        return nil, err
    }
	return &model.SlackAttachment{
		Color: "#95b7d0",
		Text: table,
		Title: "sonar",
		AuthorName: "SonarQube",
		AuthorLink: w.ServerURL,
	}, nil

}
