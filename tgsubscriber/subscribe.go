package tgsubscriber

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	SUBURL = "http://10.10.126.2:8000"
)

type ChannelJSON struct {
	Channel string `json:"channel"`
}

func SubscribeChannel(link string) error {
	hc := http.Client{}

	jsonBody := &ChannelJSON{
		Channel: link,
	}

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", SUBURL+"/subscribe", bytes.NewBuffer(jsonData))
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

	return nil
}