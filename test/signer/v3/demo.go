package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"sigs.k8s.io/cloud-provider-huaweicloud/pkg/cloudprovider/huaweicloud"
	"sigs.k8s.io/cloud-provider-huaweicloud/test/signer/v3/legacysdk/core"
)

const (
	configPathEnv = "CCMConfigPath"
	serverIDEnv   = "ServerID"
)

func readConfigFromEnv() (config *huaweicloud.CloudConfig, serverID string) {
	configPath := os.Getenv(configPathEnv)
	if len(configPath) == 0 {
		fmt.Printf("Please set config file path. e.g. 'export CCMConfigPath=/etc/provider.conf'.\n")
		return nil, ""
	}

	fileHandler, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Failed to open file: %s, error: %v\n", configPath, err)
		return nil, ""
	}
	defer fileHandler.Close()

	config, err = huaweicloud.ReadConf(fileHandler)
	if err != nil {
		fmt.Printf("Failed to parse config file: %s.\n", configPath)
		return nil, ""
	}

	serverID = os.Getenv(serverIDEnv)
	if len(serverID) == 0 {
		fmt.Printf("Please set server ID. e.g. 'export ServerID=a44af098-7548-4519-8243-a88ba3e5de4g'.\n")
		return nil, ""
	}

	return config, serverID
}

func main() {
	config, serverID := readConfigFromEnv()
	if config == nil || len(serverID) == 0 {
		return
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s:443/v1/%s/cloudservers/%s", config.Auth.ECSEndpoint, config.Auth.ProjectID, serverID)
	req, err := http.NewRequest("GET", url, ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
	if err != nil {
		return
	}
	req.Close = true

	// add the sign to request header if needed.
	sign := core.Signer{
		AccessKey: config.Auth.AccessKey,
		SecretKey: config.Auth.SecretKey,
		Region:    config.Auth.Region,
		Service:   "ec2",
	}
	req.Header.Add("X-Project-Id", config.Auth.ProjectID)
	// req.Header.Add("x-sdk-date", "20190829T122203Z")
	if err := sign.Sign(req); err != nil {
		fmt.Printf("sign failed with error: %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("request failed with error: %v\n", err)
		return
	}

	var servers huaweicloud.EcsServers
	err = huaweicloud.DecodeBody(resp, &servers)
	if err != nil {
		fmt.Printf("decode response body failed with error: %v\n", err)
		return
	}

	fmt.Printf("Congratulations! Your authentication information is suitable.\n")
}
