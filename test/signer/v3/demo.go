package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"sigs.k8s.io/cloud-provider-huaweicloud/test/signer/v3/legacysdk/core"
)

func main() {
	demoApp()
}

func demoApp() {
	// Create the HTTP request
	req, err := http.NewRequest("GET", "http://vpc.cn-north-4.myhuaweicloud.com:443/v2.0/lbaas/loadbalancers?vip_address=10.38.173.101", ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
	if err != nil {
		return
	}
	req.Close = true

	// add the sign to request header if needed.
	sign := core.Signer{
		AccessKey: "cTNDcpJDkkBM8ZP72pwd",
		SecretKey: "qpvqfcp4FcmCehcQdS1fExcXo03bGyvbaf5ugE9c",
		Region:    "cn-north-4",
		Service:   "ec2",
	}
	req.Header.Add("X-Project-Id", "88eecb5e3034473c817407e590eb3bca")
	req.Header.Add("x-sdk-date", "20190829T122203Z")
	if err := sign.Sign(req); err != nil {
		fmt.Printf("sign failed with error: %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("request failed with error: %v\n", err)
		return
	}
	fmt.Printf("response HTTP code: %v\n", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read response body failed with error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Body: %v\n", string(body))

	fmt.Printf("Congratulations! Your authentication information is suitable.\n")
}
