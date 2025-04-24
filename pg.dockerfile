FROM golang:1.24-alpine
LABEL AUTHOR="Impervguin"

RUN apk add --no-cache make

RUN mkdir /build
WORKDIR /build

COPY pkg/pg-migrations/ .
COPY ./config/migr.yaml ./config/pgmigr.yaml

RUN go mod download

ENV PATH="$PATH:/build"

RUN make build