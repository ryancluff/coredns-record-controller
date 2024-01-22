package pfsense_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
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

func (pfc *PfsenseClient) GetHostOverrides() (map[string]HostOverride, error) {
	hostOverridesByFQDN := map[string]HostOverride{}
	response, err := pfc.Call("GET", "/api/v1/services/unbound/host_override", nil)
	if err != nil {
		return hostOverridesByFQDN, err
	}
	if response["status"] != "ok" {
		err = fmt.Errorf("response status: %.0f %s", response["code"].(float64), response["status"].(string))
		return hostOverridesByFQDN, err
	}

	for _, value := range response["data"].([]interface{}) {
		valueMap := value.(map[string]interface{})
		hostOverride := HostOverride{
			Host:   valueMap["host"].(string),
			Domain: valueMap["domain"].(string),
			IP:     []string{valueMap["ip"].(string)},
			Tag:    valueMap["descr"].(string),
		}

		fqdn := fmt.Sprintf("%s.%s", hostOverride.Host, hostOverride.Domain)
		hostOverridesByFQDN[fqdn] = hostOverride
	}

	return hostOverridesByFQDN, nil
}

func (pfc *PfsenseClient) SetHostOverrides(hostOverridesByFQDN map[string]HostOverride) error {
	var hostOverrides HostOverrides
	for _, hostOverride := range hostOverridesByFQDN {
		hostOverrides = append(hostOverrides, hostOverride)
	}

	sort.Sort(hostOverrides)

	payload := map[string]interface{}{
		"host_overrides": hostOverrides,
		"apply":          true,
	}

	response, err := pfc.Call("PUT", "/api/v1/services/unbound/host_override/flush", payload)
	if err != nil {
		return err
	}
	if response["status"] != "ok" {
		err = fmt.Errorf("response status: %.0f %s", response["code"].(float64), response["status"].(string))
		return err
	}

	return nil
}
