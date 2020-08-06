/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"sigs.k8s.io/cloud-provider-huaweicloud/pkg/legacy/apigwv3/core"
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
		return
	}

}
