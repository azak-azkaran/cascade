FROM golang:alpine AS build-env
RUN apk --no-cache add git make
RUN go get  github.com/azak-azkaran/cascade
WORKDIR /go/src/github.com/azak-azkaran/cascade
RUN make install

FROM alpine
# less priviledge user, the id should map the user the downloaded files belongs to
RUN apk --no-cache add shadow && \
        groupadd -r dummy && \
        useradd -r -g dummy dummy -u 1000

COPY --from=build-env /go/bin/cascade /cascade

CMD ["./cascade", "-c", "/tmp/config.yaml"]
