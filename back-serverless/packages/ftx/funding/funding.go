package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Request takes in the user's input for the future they want and if the type is a GET or PUT.
type Request struct {
	Future    string `json:"future"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// Response returns back the http code, type of data, and the presigned url to the user.
type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       interface{}       `json:"body,omitempty"`
}

var (
	ErrNoFuture    = errors.New("no future provided")
	ErrNoStartDate = errors.New("no start date provided")
	ErrNoEndDate   = errors.New("no end date provided")
	ErrFTXResponse = errors.New("ftx returned something else than a 200")
)

type FTXResponse struct {
	Success bool `json:"success"`
	Result  []struct {
		Future string    `json:"future"`
		Rate   float64   `json:"rate"`
		Time   time.Time `json:"time"`
	} `json:"result"`
}

func Main(in Request) *Response {
	if in.Future == "" {
		return &Response{StatusCode: http.StatusBadRequest, Body: ErrNoFuture.Error()}
	}

	startTime, err := time.Parse("2006-01-02", in.StartTime)
	if err != nil {
		return &Response{StatusCode: http.StatusBadRequest, Body: ErrNoStartDate.Error()}
	}

	endTime, err := time.Parse("2006-01-02", in.EndTime)
	if err != nil {
		return &Response{StatusCode: http.StatusBadRequest, Body: ErrNoEndDate.Error()}
	}

	url := fmt.Sprintf("https://ftx.com/api/funding_rates?future=%s&start_time=%d&end_time=%d", in.Future, startTime.Unix(), endTime.Unix())

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &Response{StatusCode: http.StatusBadRequest, Body: err.Error()}
	}
	res, err := client.Do(req)
	if err != nil {
		return &Response{StatusCode: http.StatusBadRequest, Body: err.Error()}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return &Response{StatusCode: http.StatusBadRequest, Body: ErrFTXResponse.Error()}
	}

	var j FTXResponse
	err = json.NewDecoder(res.Body).Decode(&j)
	if err != nil {
		panic(err)
	}

	// json headers
	return &Response{StatusCode: http.StatusOK, Headers: map[string]string{
		"Content-Type": "application/json",
	}, Body: j}
}
