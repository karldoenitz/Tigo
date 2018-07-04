// Copyright 2018 The Tigo Authors. All rights reserved.
package WebFramework

import (
	"net/http"
	"fmt"
)

type Application struct {
	IPAddress  string
	Port       string
	UrlPattern UrlPattern
}

func (application *Application)Run() {
	application.UrlPattern.Init()
	address := fmt.Sprintf("%s:%s", application.IPAddress, application.Port)
	httpServerErr := http.ListenAndServe(address, nil)
	if httpServerErr != nil {
	}
}
