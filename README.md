# cascade
[![Build Status](https://travis-ci.org/azak-azkaran/cascade.svg?branch=master)](https://travis-ci.org/azak-azkaran/cascade)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=azak-azkaran_cascade&metric=alert_status)](https://sonarcloud.io/dashboard?id=azak-azkaran_cascade)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=azak-azkaran_cascade&metric=coverage)](https://sonarcloud.io/dashboard?id=azak-azkaran_cascade)
[![Coverage Status](https://coveralls.io/repos/github/azak-azkaran/cascade/badge.svg)](https://coveralls.io/github/azak-azkaran/cascade)

Go proxy which can switch between Direct mode and Cascade mode
Switch is done according to health check.

## Installation
Just download the current release
For running as daemon checkout the section in this readme

### Install from source
checkout the source code and run 
```
make install
```

## Configuration
Configuration can be done by file or command line arguments

* __password__ : Password for authentication to a forward proxy
* __host__ : Address of a forward proxy
* __user__ : Username for authentication to a forward proxy
* __port__ : Port on which to run the proxy
* __health__ : Address which is used for health check if available go to direct mode (default: https://www.google.de )
* __health-time__ : Duration between health checks (default: 30 Seconds )
* __host-list__ : Comma Separated List of Host for which Proxy Redirect is used in Cascade Mode
* __config__ : Path to configuration yaml file. If set all other command line parameters will be ignored
* __version__: Just shows the current version
* __vault-addr__: address to the hashi corp vault
* __vault-token__: token for the hashi corp vault

### Health Check
A temporary client tries to connect to a certain address regularly.
The Cascade mode is active if health check fails.

### Direct Mode
Normal http internet Proxy Mode.

### Cascade Mode
Cascade Mode means that this proxy stands between the client and another Proxy.
Basic Authentication can be enabled for Cascade Mode

### Proxy redirect for Hosts

If in cascade mode, different Proxy redirects can be added by adding a Comma seperated list. These redirects changes HTTPS and HTTP Request according to the setup rules.
The Requests can be send to another Proxy or directly.

Direct Configuration:

* __DIRECT Connection__: eclipse
* __DIRECT Connection__: azure->
* __REDIRECT Connection to other Proxy__: google->test:8888

### HashiCorp Vault Configuration
Cascade can also get the configuration from a HashiCorp Vault.
The secret store has to be version 2 and a vault configuration will overwrite a file configuration.
The following secrets have to be saved in store:

* username
* password
* host 

The following secrets are optional:

* port (default: 8888) 
* host-list (default: "")
* health (default: www.google.de) 
* log (default: WARNING)
* health-time (default: 5s)
* disableAutoChangeMode  (default: false)


## REST Interface

Cascade comes with a REST Interface which can be used to control the application.
Currently the following REST Endpoints are available:

* __/config__ : is used to return the current configuration
* __/getOnlineCheck__ : is used to get if the check is used to check for up 
* __/getAutoMode__ : is used to get if automatically switching between modes is activated
* __/addRedirect__ : adds another redirect rule
* __/setOnlineCheck__ : used to configure check is used to check for a website being online
* __/setAutoMode__ : used to disable automatically switching between modes
* __/setCascadeMode__ : used to force a certain mode

### Curl Examples
These examples use curl to use the REST Endpoints

#### Config
To get the current configuration
```
curl -D- localhost/config
```

#### SetCascadeMode
To set the Mode by hand to DirectMode:
```
curl -D- --request POST \
  --data '{"cascadeMode":false}' \
  localhost/setCascadeMode
```
To set the Mode by hand to Cascade mode:
```
curl -D- --request POST \
  --data '{"cascadeMode":true}' \ 
  localhost/setCascadeMode
```

#### SetAutoMode
To disable the automatically switch between modes:
```
curl -D- --request POST \
  --data '{"autoChangeMode":false}' \
  localhost/setAutoMode
```
To enable the automatically switch between modes:
```
curl -D- --request POST \
  --data '{"autoChangeMode":true}' \
  localhost/setAutoMode
```

**will overwrite the setting from file,vault or command line**

## Systemd

If you want to use the provided service configuration, the program has to be moved to 
```
/usr/local/bin
```
The configuration has to be moved to the following folder for Ubuntu:

```
/etc/systemd/system/
```

Afterwards, systemd has to be restarted as follows:
```
systemctl daemon-reload
systemctl start cascade
```

The logs can be viewed by using:
```
journalctl -f -u cascade
```

