FROM golang:1.24-alpine
LABEL AUTHOR="Impervguin"

RUN apk add --no-cache make curl libstdc++ libgcc

RUN mkdir /logs
RUN mkdir /build
WORKDIR /build

COPY go.* ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go get -tool github.com/a-h/templ/cmd/templ@latest

COPY ./cmd/api/ ./cmd/api/
COPY ./internal/ ./internal/

RUN swag init --parseInternal --parseDependency --parseDepth 10 -g ./cmd/api/main.go -o ./cmd/docs
RUN go tool templ generate -path ./internal/view
RUN go build -o ./api.bin ./cmd/api

CMD ["./api.bin"]

