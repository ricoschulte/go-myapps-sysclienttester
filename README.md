# go-myapps-sysclienttester

A tool for adding dummy devices to the device app of an innovaphone myApps environment via the sysclient protocol.



Useful for testing during app development, debugging the configuration within the device app or similar.



- can connect to the websocket sysclient url of device apps like to any other device

- can emulate one or more devices simultaneously

- either serves a static web page or can serve files from a local directory, which can then be accessed through the device app's user interface

- creates a logstream on stdout with the configuration that the device app sends to the emulated devices, such as STUN/TURN configuration, NTP, or any other configuration that the device app can run.

# Build

To build a binary version of the app use the go compiler as usual:

``` BASH
go build -o go-myapps-sysclienttester .
```

As with any other go code, crosscompile it to other platforms than the host you are compiling on:

``` BASH
# Linux amd64
GOOS="linux" GOARCH="amd64" go build -o go-myapps-sysclienttester .

# Raspberry Pi 3
GOOS="linux" GOARCH="arm" go build -o go-myapps-sysclienttester .

# Raspberry Pi 4 
GOOS="linux" GOARCH="arm64" go build -o go-myapps-sysclienttester .

# Windows exe
GOOS="windows" GOARCH="amd64" go build -o go-myapps-sysclienttester.exe .
```

# Usage on the shell

``` BASH
./go-myapps-sysclienttester \
  -secretkey yoursecretkey \
  -sysclient wss://apps.company.com/company.com/devices/sysclients \
  -idfile identities.json \
  -staticdir ./static/ \
  -sessiondir ./sessions/ \
  -loglevel trace
```

## Command line options

Usage of /go-myapps-sysclienttester:
```
  -idfile string
        path to a JSON file with device identities (default "identities.json")
  -loglevel string
        log level to use. (default "info")
  -secretkey string
        secretkey used to encrypt local session files
  -sessiondir string
        path to a local folder to save the sessions
  -skipverifytls
        disabled verification of TLS certificates
  -staticdir string
        path to a local folder to serve static files under <device-passthrough>/static
  -sysclient string
        the sysclient url of the devices app to connect to (wss://apps.company.com/company.com/devices/sysclients)
```

## Format identity file

The file defines a key `ids`. In that key is a list of identities to emulate.

``` JSON
{
  "ids": [
    {
      "id": "f049033123450",
      "product": "IP990",
      "version": "IP990",
      "fwBuild": "900001",
      "bcBuild": "900001",
      "major": "v0.0.0",
      "fw": "firmware.bin",
      "bc": "bootcode.bin",
      "mini": false,
      "pbxActive": false,
      "other": false,
      "platform": {
        "type": "PHONE",
        "fxs": false
      },
      "ethIfs": [
        {
          "if": "ETH0",
          "ipv4": "172.16.0.141",
          "ipv6": "2002:9110:9d07:0:290:33ff:fe46:af0"
        },
        {
          "if": "ETH1",
          "ipv4": "172.16.1.141",
          "ipv6": "2002:9111:9d07:0:290:33ff:fe46:af0"
        },
        {
          "if": "ETH2",
          "ipv4": "172.16.2.141",
          "ipv6": "2002:9112:9d07:0:290:33ff:fe46:af0"
        }
      ]
    }
    
    // insert more identities here
    // you can serve as many dummy devices as you like 
    // at the same time
    
  ]
}
```
The description of the individual keys can be found at the [SDK Sysclient protocol docs](https://sdk.innovaphone.com/13r3/doc/protocol/sysclient.htm#AMTIdentify)

# Docker

To build a Docker container with the app:

``` BASH
# build
docker build -t your-container:latest .

# run
docker run \
  --rm -it \
  -v ${PWD}/identities.json:/data/identities.json \
  -v ${PWD}/static:/data/static \
  -v ${PWD}/sessions:/data/sessions \
  your-container:latest \
  -secretkey yoursecretkey \
  -sysclient wss://apps.company.com/company.com/devices/sysclients \
  -idfile /data/identities.json \
  -staticdir /data/static/ \
  -sessiondir /data/sessions/ \
  -loglevel trace
```

or use a `docker-compose.yml`:

``` YAML
version: "3"
services:
  go-myapps-sysclienttester:
    build:
      context: .
      dockerfile: Dockerfile
    restart: always
    volumes:
      - ${PWD}/identities.json:/data/identities.json:ro
      - ${PWD}/static:/data/static:ro
      - ${PWD}/sessions:/data/sessions:rw
    command: |
      -secretkey yoursecretkey 
      -sysclient wss://apps.company.com/company.com/devices/sysclients
      -idfile /data/identities.json
      -staticdir /data/static/
      -sessiondir /data/sessions/
      -loglevel trace     
```

``` BASH
docker-compose up -d --build
```

All commandline options can be used 

## About ©

[myApps](https://www.innovaphone.com/en/myapps/what-is-myapps.html) is a product of [innovaphone AG](https://www.innovaphone.com).

Documentation of the API used in this client can be found at [ innovaphone App SDK](https://sdk.innovaphone.com/).