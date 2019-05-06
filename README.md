# cascade
[![Build Status](https://travis-ci.org/azak-azkaran/cascade.svg?branch=master)](https://travis-ci.org/azak-azkaran/cascade)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=azak-azkaran_cascade&metric=alert_status)](https://sonarcloud.io/dashboard?id=azak-azkaran_cascade)
[![Coverage Status](https://coveralls.io/repos/github/azak-azkaran/cascade/badge.svg?branch=master)](https://coveralls.io/github/azak-azkaran/cascade?branch=master)

golang proxy which can switch between Direct mode and Cascade mode
Swich is done according to health check.

## Health Check
A temporary client tries to connect to a certain address regurlary.
The Cascade mode is active if health check fails.

## Direct Mode
Normal http internet Proxy Mode.

## Cascade Mode
Cascade Mode means that this proxy stands between the client and another Proxy.
Basic Auth can be enabled for Cascade Mode
