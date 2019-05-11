# cascade
[![Build Status](https://travis-ci.org/azak-azkaran/cascade.svg?branch=master)](https://travis-ci.org/azak-azkaran/cascade)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=azak-azkaran_cascade&metric=alert_status)](https://sonarcloud.io/dashboard?id=azak-azkaran_cascade)
[![Coverage Status](https://coveralls.io/repos/github/azak-azkaran/cascade/badge.svg?branch=master)](https://coveralls.io/github/azak-azkaran/cascade?branch=master)

golang proxy which can switch between Direct mode and Cascade mode
Swich is done according to health check.

## Configuration
Configuration can be done by file or command line argumments.

* __password__ : Password for authentication to a forward proxy
* __host__ : Address of a forward proxy
* __user__ : Username for authentication to a forward proxy
* __port__ : Port on which to run the proxy
* __health__ : Address which is used for health check if available go to direct mode (default: https://www.google.de )
* __health-time__ : Duration between health checks (default: 30 Seconds )
* __host-list__ : Comma Separated List of Host for which DirectMode is used in Cascade Mode
* __config__ : Path to config yaml file. If set all other command line parameters will be ignored

## Health Check
A temporary client tries to connect to a certain address regurlary.
The Cascade mode is active if health check fails.

## Direct Mode
Normal http internet Proxy Mode.

## Cascade Mode
Cascade Mode means that this proxy stands between the client and another Proxy.
Basic Auth can be enabled for Cascade Mode

### Skip Cascade for Hosts

If in cascade mode, requests to the provided hosts will be send diretly without sending it to the forward proxy
