package api

import (
	"fmt"

	"github.com/imroc/req"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func printError(resp *req.Resp) {
	errorMessage := ErrorMessage{}
	resp.ToJSON(&errorMessage)
	fmt.Println(errorMessage)
}
