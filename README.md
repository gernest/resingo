# resingo [![GoDoc](https://godoc.org/github.com/gernest/resingo?status.svg)](https://godoc.org/github.com/gernest/resingo) [![Go Report Card](https://goreportcard.com/badge/github.com/gernest/resingo)](https://goreportcard.com/report/github.com/gernest/resingo)

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
 - [ ] Remove/Delete device
 - [ ] Identify device
 - [x] Rename device
 - [ ] Note a device
 - [x] Generate device UUID
 - [x] Register device
 - [ ] Enable device url
 - [ ] Disable device url
 - [ ] Move device
 - [ ] Check device status

- Environment
 - Device
  - [ ] Get all device environment variables
  - [ ] Create a device environment variable
  - [ ] Update device environment variable
  - [ ] Remove device environment variable
 - Application
  - [ ] Get all application environment variables
  - [ ] Create application environment variable
  - [ ] Update application environment variable
  - [ ] Remove application environment variable

- Keys
 - [ ] Get all ssh keys
 - [ ] Get a dingle ssh key
 - [ ] Remove ssh key
 - [ ] Create ssh key

- Os
  - [ ] Download Os Image

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
