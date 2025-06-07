package main

import (
	"bytes"
	"encoding/json"
	"log"
	"config"
	"internal/news"
	"os"
	"time"
	"github.com/joho/godotenv"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type News struct {
	Text   string `json:"text"`
	Time   string `json:"time"`
	Source string `json:"source"`
	URL    string `json:"url"`
}

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
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	db := news.ConnectDB(cfg)
	defer db.Close()

	repo := news.NewRepository(db)

	jsons := []string{
		`{
			"text": "SpaceX успешно провела испытание системы аварийного спасения экипажа...",
			"time": "2025-06-07T06:50:00Z",
			"source": "BBC Science",
			"url": "https://bbc.com/news/science/spacex-starship-crew-test"
		}`,
		`{
			"text": "Ракета Starship от SpaceX прошла ключевой тест...",
			"time": "2025-06-07T07:05:00Z",
			"source": "РИА Новости",
			"url": "https://ria.ru/20250607/spacex-starship-test-1850000000.html"
		}`,
	}

	for _, js := range jsons {
		var n News
		_ = json.Unmarshal([]byte(js), &n)
		parsedTime, _ := time.Parse(time.RFC3339, n.Time)
		err := repo.Insert(n.Text, parsedTime, n.Source, n.URL)
		if err != nil {
			log.Println("❌ Ошибка вставки:", err)
		}
	}

	log.Println("✅ Новости добавлены")

	go func() {
		err := sendMessage("Привет!")
		if err != nil {
			fmt.Printf("%++v\n", err)
			return
		}

		time.Sleep(10 * time.Second)

		answer, err := getResponse()
		if err != nil {
			fmt.Printf("%++v\n", err)
			return
		}

		fmt.Printf("%s\n", answer)

		err = clearContext()
		if err != nil {
			fmt.Printf("%++v\n", err)
			return
		}
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

type ResponseRequestJSON struct {
	OperatingSystemCode int    `json:"operatingSystemCode"`
	Dialogidentifier    string `json:"Dialogidentifier"`
	ApiKey              string `json:"apiKey"`
}

func getResponse() (string, error) {
	hc := http.Client{}

	jsonBody := &ResponseRequestJSON{
		OperatingSystemCode: 12,
		ApiKey:              APIKEY,
		Dialogidentifier:    " " + USERNAME + "_1",
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

func clearContext() error {
	hc := http.Client{}

	jsonBody := &ResponseRequestJSON{
		OperatingSystemCode: 12,
		ApiKey:              APIKEY,
		Dialogidentifier:    " " + USERNAME + "_1",
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

// func enableCors(writer *http.ResponseWriter) {
// 	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
// }
