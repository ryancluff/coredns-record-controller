package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PfsenseClient struct {
	host   string
	auth   string
	client *http.Client
}

func NewPfsenseClient(host string, clientID string, clientToken string) (PfsenseClient, error) {
	pfc := PfsenseClient{
		host:   host,
		auth:   clientID + " " + clientToken,
		client: &http.Client{},
	}

	response, err := pfc.Call("GET", "/api/v1/services/unbound/host_override", nil)
	if err != nil {
		err = fmt.Errorf("pfsense connection test failed; %w", err)
		return pfc, err
	}
	if response["status"] != "ok" {
		err = fmt.Errorf("pfsense connection test failed; response status: %.0f %s", response["code"].(float64), response["status"].(string))
		return pfc, err
	}

	return pfc, nil
}

func (pfc *PfsenseClient) Call(method string, path string, payload map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	url := "https://" + pfc.host + path

	payloadString, err := json.Marshal(payload)
	payloadReader := bytes.NewReader(payloadString)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest(method, url, payloadReader)
	if err != nil {
		return result, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", pfc.auth)

	res, err := pfc.client.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}

	return result, nil
}
