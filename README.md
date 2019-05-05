# cascade
golang proxy which can switch between Direct mode and Cascade mode
Swich is done according to health check.

## Health Check
A temporary client tries to connect to a certain address regurlary.
The Cascade mode is active if health check succeeds.

__TODO__ This will be changed in the future.

## Direct Mode
Normal http internet Proxy Mode.

## Cascade Mode
Cascade Mode means that this proxy stands between the client and another Proxy.
Basic Auth can be enabled for Cascade Mode
