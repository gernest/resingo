# resingo [![GoDoc](https://godoc.org/github.com/gernest/resingo?status.svg)](https://godoc.org/github.com/gernest/resingo) [![Go Report Card](https://goreportcard.com/badge/github.com/gernest/resingo)](https://goreportcard.com/report/github.com/gernest/resingo)[![Build Status](https://travis-ci.org/gernest/resingo.svg?branch=master)](https://travis-ci.org/gernest/resingo)[![Coverage Status](https://coveralls.io/repos/github/gernest/resingo/badge.svg?branch=master)](https://coveralls.io/github/gernest/resingo?branch=master)

The unofficial golang sdk for resin.io

## what is resin and what is resingo?
[Resin](https://resin.io/) is a service that brings the benefits of linux
containers to IoT. It enables iterative development, deployment and management
of devices at scale. Please go to their website for more information.

Resingo, is a standard development kit for resin, that uses the Go(Golang)
programming language. This library, provides nuts and bolts necessary for
developers to interact with resin service using the Go programming language.

This library is full featured, well tested and well designed. If you find
anything that is missing please [open an issssue](https://github.com/gernest/resingo/issues)



This is the laundry list of things you can do.

- Applications
 - [x] Get all applications
 - [x] Get application by name
 - [x] Get application by application id
 - [x] Create Application
 - [x] Delete application
 - [x] Generate application API key
- Devices
 - [x] Get all devices
 - [x] Get all devices for a given device name
 - [x] Get device by UUID
 - [x] Get device by name
 - [x] Get application name by device uuid
 - [x] Check if the device is online or not
 - [x] Get local IP address of the device
 - [x] Remove/Delete device
 - [ ] Identify device
 - [x] Rename device
 - [x] Note a device
 - [x] Generate device UUID
 - [x] Register device
 - [x] Enable device url
 - [x] Disable device url
 - [x] Move device
 - [x] Check device status
 - [x] Identify device by blinkig

- Environment
 - Device
  - [x] Get all device environment variables
  - [x] Create a device environment variable
  - [x] Update device environment variable
  - [x] Remove/Delete device environment variable
 - Application
  - [x] Get all application environment variables
  - [x] Create application environment variable
  - [x] Update application environment variable
  - [x] Remove application environment variable

- Keys
 - [x] Get all ssh keys
 - [x] Get a dingle ssh key
 - [x] Remove ssh key
 - [x] Create ssh key

- Os
  - [ ] Download Os Image

- Config
 - [x] Get all configurations

- Logs
 - [x] Subscribe to device logs
 - [ ] Retrieve historical logs

 - Supervisor
  - [x] Reboot

 # Introduction

 ## Installation

 ```bash
 go get github.com/gernest/resingo
 ```

## Design philosophy

#### Naming convention
The library covers different componets of the resin service namely  `device`,
`application`, `os` ,`environment` and `keys`.

Due to the nature of this lirary to use functions rather than attach methods to
the relevant structs. We use a simple strategy of naming functions that operates
for the different components. This is by adding Prefix that represent the resin
component.

The following is the Prefix table that shows all the prefix used to name the
functions. For example `DevGetAll` is a function that retrives all devices and
`AppGetAll` is a function that retrieves all App;ications.

component   | function prefix
------------|----------------
Device      | Dev
Application | App
EnvironMent | Env
Os          | Os
Keys        | Key

#### Functions over methods
This library favors functions to provide functionality that interact with resin.
All this functions accepts `Context` as the first argument, and can optionally
accepts other arguments which are function specific( Note that this `Context` is
defined in this libary, don't mistook it for the `context.Context`).

The reason behind this is to provide composability and making the codebase
clean. Also it is much easier for testing.


# Usage

```go

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gernest/resingo"
)

func main() {

	// We need the context from which we will be talking to the resin API
	ctx := &resingo.Context{
		Client: &http.Client{},
		Config: &resingo.Config{
			Username: "Your resing username",
			Password: "your resubn password",
		},
	}

	/// There are two ways to authenticate the context ctx.

	// Authenticate with credentials i.e Username and password
	err := resingo.Login(ctx, resingo.Credentials)
	if err != nil {
		log.Fatal(err)
	}

	// Or using authentication token, which you can easily find on your resin
	// dashboard
	err = resingo.Login(ctx, resingo.AuthToken, "Tour authentication token goes here")
	if err != nil {
		log.Fatal(err)
	}

	// Now the ctx is authenticated you can pass it as the first arument to any
	// resingo API function.
	//
	// The ctx is safe for concurrent use

	// Get All your applications
	apps, err := resingo.AppGetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range apps {
		fmt.Println(*a)
	}
}
```



# Contributing

This requires go1.7+

Running tests will require a valid resin account . You need to set the following
environment variables before running the `make` command.

```bash
export RESINTEST_EMAIL=
export RESINTEST_PASSWORD=
export RESINTEST_USERNAME=
export RESINTEST_REALDEVICE_UUID=
```

The names are self explanatory. To avoid typing them all the time, you can write
them into a a file named `.env` which stays at the root of this project, the test
script will automatically source it for you.

All contributions are welcome.

# Author

twitter [@gernesti](https://twitter.com/gernesti)

# Licence
MIT
