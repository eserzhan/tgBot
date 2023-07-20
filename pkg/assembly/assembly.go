package assembly

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
	TRANSCRIPT_URL = "https://api.assemblyai.com/v2/transcript"
	statusCompleted = "completed"
	statusProcessing = "processing"
	contentHeader = "content-type"
	authHeader = "authorization"
	headerValue = "application/json"
)

// Client is a Assembly API client
type Client struct {
	client      *http.Client
	apiKey string
}


func NewClient(key string) (*Client, error) {
	if key == "" {
		return nil, errors.New("api key is empty")
	}

	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		apiKey: key,
	}, nil
}

func(c *Client) Transcription(url string) (string, error) {
	values := map[string]string{"audio_url": url}
	jsonData, err := json.Marshal(values)
	if err != nil {
		return "", err 
	}
	req, _ := http.NewRequest(http.MethodPost, TRANSCRIPT_URL, bytes.NewBuffer(jsonData))
	req.Header.Set(contentHeader, headerValue)
	req.Header.Set(authHeader, c.apiKey)
	res, err := c.client.Do(req)
	if err != nil {
		return "", err 
	}
	defer res.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)


	return result["id"].(string), nil
}

func(c *Client) TranscribedText(s string) (string, error) {
	POLLING_URL := TRANSCRIPT_URL + "/"	+ s
	req, _ := http.NewRequest(http.MethodGet, POLLING_URL, nil)
	req.Header.Set(contentHeader, headerValue)
	req.Header.Set(authHeader, c.apiKey)

	for{
	res, err := c.client.Do(req)
	if err != nil {
		return "", err 
	}

	defer res.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return "", err 
	}

	if result["status"] == statusCompleted {
		return result["text"].(string), nil 
	}else if result["status"] == statusProcessing{
		time.Sleep(1 * time.Second)
		continue 
	}else{
		return "", errors.New(result["status"].(string))
	}
}
}