FROM golang:1.24-alpine
LABEL AUTHOR="Impervguin"

RUN apk add --no-cache make

RUN mkdir /logs
RUN mkdir /build
WORKDIR /build


COPY go.* .
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest


COPY ./cmd/api/ ./cmd/api/
COPY ./internal/ ./internal/
COPY ./config/*.yaml ./config/
COPY ./makefile .

CMD ["make", "api"]

