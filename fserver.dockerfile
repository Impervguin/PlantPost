FROM golang:1.24-alpine
LABEL AUTHOR="Impervguin"

RUN apk add --no-cache make

RUN mkdir /build
WORKDIR /build

COPY pkg/fserver/ .
COPY ./config/fserver.yaml ./config/fserver.yaml

RUN go mod download

RUN make build