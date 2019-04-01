package analysis

import (
	"fmt"

	"github.com/imroc/req"
)

type SentimentDatum struct {
	Value     string `json:"value"`
	ValueName string `json:"value_classification"`
	Timestamp string `json:"timestamp"`
	TimeUntil string `json:"time_until_update"`
}
type Sentiment struct {
	Name     string        `json:"name"`
	Data     SentimentData `json:"data"`
	Metadata interface{}   `json:"metadata"`
}
type SentimentData []SentimentDatum

func GetSentimentIndex() *SentimentDatum {
	resp, err := req.Get("https://api.alternative.me/fng/")
	sentiment := Sentiment{}
	if err != nil || resp.Response().StatusCode != 200 {
		fmt.Println(err)
		return nil
	}
	resp.ToJSON(&sentiment)
	return &sentiment.Data[0]
}
