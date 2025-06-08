package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ResponseStatusJSON struct {
	IsSuccess   bool   `json:"isSuccess"`
	Description string `json:"description"`
}

type ResponseDataJson struct {
	UserCrocCode     string `json:"userCrocCode"`
	DialogIdentifier string `json:"dialogIdentifier"`
	LastMessage      string `json:"lastMessage"`
	LastResponseTime string `json:"lastResponseTime"`
}

type ResponseJSON struct {
	Status ResponseStatusJSON `json:"status"`
	Data   ResponseDataJson   `json:"data"`
}

type RequestJSON struct {
	OperatingSystemCode int    `json:"operatingSystemCode"`
	ApiKey              string `json:"apiKey"`
	UserDomainName      string `json:"userDomainName"`
	DialogIdentifier    string `json:"dialogIdentifier"`
	AiModelCode         int    `json:"aiModelCode"`
	Message             string `json:"Message"`
}

type ResponseRequestJSON struct {
	OperatingSystemCode int    `json:"operatingSystemCode"`
	Dialogidentifier    string `json:"Dialogidentifier"`
	ApiKey              string `json:"apiKey"`
}

func sendMessage(message string, id string) error {
	hc := http.Client{}

	jsonBody := &RequestJSON{
		OperatingSystemCode: 12,
		ApiKey:              APIKEY,
		UserDomainName:      USERNAME,
		DialogIdentifier:    " " + USERNAME + "_" + id,
		AiModelCode:         1,
		Message:             message,
	}

	jsonData, err := json.MarshalIndent(
		jsonBody,
		"",
		"\t",
	)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", APIURL+"/PostNewRequest", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	var data ResponseStatusJSON
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	if !data.IsSuccess {
		return errors.New(data.Description)
	}

	return nil
}

func getResponse(id string) (string, error) {
	hc := http.Client{}

	jsonBody := &ResponseRequestJSON{
		OperatingSystemCode: 12,
		ApiKey:              APIKEY,
		Dialogidentifier:    " " + USERNAME + "_" + id,
	}

	jsonData, err := json.MarshalIndent(
		jsonBody,
		"",
		"\t",
	)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", APIURL+"/GetNewResponse", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	var data ResponseJSON
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	if !data.Status.IsSuccess {
		return "", errors.New(data.Status.Description)
	}

	return data.Data.LastMessage, nil
}

func clearContext(id string) error {
	hc := http.Client{}

	jsonBody := &ResponseRequestJSON{
		OperatingSystemCode: 12,
		ApiKey:              APIKEY,
		Dialogidentifier:    " " + USERNAME + "_" + id,
	}

	jsonData, err := json.MarshalIndent(
		jsonBody,
		"",
		"\t",
	)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", APIURL+"/CompleteSession", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	var data ResponseStatusJSON
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	if !data.IsSuccess {
		return errors.New(data.Description)
	}

	return nil
}

func ObtainRequest(message string) (string, error) {
	id := strconv.FormatInt(time.Now().Unix(), 10)
	err := sendMessage(message, id)
	if err != nil {
		fmt.Printf("%++v\n", err)
		return "", err
	}

	time.Sleep(10 * time.Second)

	result, err := getResponse(id)
	if err != nil {
		fmt.Printf("%++v\n", err)
		return "", err
	}

	err = clearContext(id)
	if err != nil {
		fmt.Printf("%++v\n", err)
		return "", err
	}

	return result, nil
}