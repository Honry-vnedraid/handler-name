package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func main() {
	go func() {
		err := sendMessage("Привет!")
		if err != nil {
			fmt.Printf("%++v\n", err)
			return
		}

		fmt.Printf("Hi\n")

		// time.Sleep(10 * time.Second)

		// answer, err := getResponse()
		// if err != nil {
		// 	fmt.Printf("%++v\n", err)
		// 	return
		// }

		// fmt.Printf("Hi\n")

		// fmt.Printf("%s\n", answer)

		// err = clearContext()
		// if err != nil {
		// 	fmt.Printf("%++v\n", err)
		// 	return
		// }
		// fmt.Printf("Hi\n")
	}()
	select {}
}

type RequestJSON struct {
	OperatingSystemCode int    `json:"operatingSystemCode"`
	ApiKey              string `json:"apiKey"`
	UserDomainName      string `json:"userDomainName"`
	DialogIdentifier    string `json:"dialogIdentifier"`
	AiModelCode         int    `json:"aiModelCode"`
	Message             string `json:"Message"`
}

func sendMessage(message string) error {
	hc := http.Client{}

	jsonBody := &RequestJSON{
		OperatingSystemCode: 12,
		ApiKey:              APIKEY,
		UserDomainName:      USERNAME,
		DialogIdentifier:    " " + USERNAME + "_1",
		AiModelCode:         1,
		Message:             message,
	}

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}

	fmt.Println(jsonBody)

	req, err := http.NewRequest("POST", APIURL+"/PostNewRequest", bytes.NewBuffer(jsonData))
	fmt.Printf("!%++v!\n", req.Body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

func getResponse() (string, error) {
	hc := http.Client{}

	form := url.Values{}
	form.Add("operatingSystemCode", "12")
	form.Add("Dialogidentifier", " Teamq7g7P8XMj52Q_1")
	form.Add("apiKey", "OrVrQoQ6T43vk0McGmHOsdvvTiX446RJ")
	req, err := http.NewRequest("POST", APIURL+"/GetNewResponse", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

func clearContext() error {
	hc := http.Client{}

	form := url.Values{}
	form.Add("operatingSystemCode", "12")
	form.Add("Dialogidentifier", " Teamq7g7P8XMj52Q_1")
	form.Add("apiKey", "OrVrQoQ6T43vk0McGmHOsdvvTiX446RJ")
	req, err := http.NewRequest("POST", APIURL+"/CompleteSession", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

// func enableCors(writer *http.ResponseWriter) {
// 	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
// }
