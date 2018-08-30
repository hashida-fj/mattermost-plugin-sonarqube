package main

import (
	// "strings"
    "bytes"
    "text/template"
    "errors"

	"github.com/mattermost/mattermost-server/model"
)

type Condition struct {
    ErrorThreshold string `json:"errorThreshold"`
    Metric         string `json:"metric"`
    OnLeakPeriod   bool   `json:"onLeakPeriod"`
    Operator       string `json:"operator"`
    Status         string `json:"status"`
    Value          string `json:"value,omitempty"`
}

type WebhookResponse struct {
	AnalysedAt string `json:"analysedAt"`
	Project    struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"project"`
	Properties struct {
	} `json:"properties"`
	QualityGate struct {
        Conditions []Condition
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"qualityGate"`
	ServerURL string `json:"serverUrl"`
	Status    string `json:"status"`
	TaskID    string `json:"taskId"`
}

type convRule struct {
    metric string
    metricjp string
    unit string
    rating bool
}

func (w *WebhookResponse) findMetric(metric string) (*Condition, error) {
    for _, cond := range w.QualityGate.Conditions {
        if cond.Metric == metric {
            return &cond, nil
        }
    }

    return nil, errors.New("unkonw metric : metric")
}

func (w *WebhookResponse) buildTable() (string, error) {
    var rows bytes.Buffer

    number2rate := map[string]string{"1": "A", "2": "B", "3": "C", "4": "D", "5": "E"}
    status2icon := map[string]string{"OK": ":white_check_mark: ", "ERROR": ":x: "}

    sonarMetricsArray := []convRule {
        {metric:"new_maintainability_rating", metricjp:"メンテナンス性", unit:"", rating:true},
        {metric:"new_reliability_rating", metricjp:"信頼性", unit:"", rating:true},
        {metric:"new_security_rating", metricjp:"セキュリティ", unit:"", rating:true},
        {metric:"new_coverage", metricjp:"カバレッジ", unit:" %", rating:false, },
        {metric:"new_duplicated_lines_density", metricjp:"コード重複率", unit:" %", rating:false},
    }

    header := "|品質指標|現在の値|閾値|\n|:-|:-|:-|\n"
    tmplstr := "|{{.Status}}{{.Metric}}|{{.Value}}|{{.ErrorThreshold}}|\n"
    rowtmpl, err:= template.New("row").Parse(tmplstr)
    if err != nil { return "", err}

    rows.WriteString(header)

    for _, rule := range sonarMetricsArray{
        cond, err:= w.findMetric(rule.metric)

        if err != nil  {
            rows.WriteString(rule.metric + "is not found")
            continue
        }

        icon := status2icon[cond.Status]
        metricjp := "[" + rule.metricjp + "](" + "http://hashida.cs.flab.fujitsu.co.jp:19000/component_measures?id=" + w.Project.Key + "&metric=" + rule.metric+ ")"
        var value string
        var threshold string
        if (rule.rating) {
            value = number2rate[cond.Value]
            threshold = number2rate[cond.ErrorThreshold]
        } else {
            value = cond.Value[:5] + rule.unit
            threshold = cond.ErrorThreshold + rule.unit
        }

        convertedCond := &Condition {
            Value:value,
            Metric:metricjp,
            Status:icon,
            ErrorThreshold:threshold,
        }

        if err := rowtmpl.Execute(&rows, convertedCond ); err != nil {
            rows.WriteString("template excution failed")
            continue
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
		AuthorName: "SonarQube: " + w.Project.Name,
		AuthorLink: "http://hashida.cs.flab.fujitsu.co.jp:19000/dashboard?id=" + w.Project.Key,
	}, nil

}
