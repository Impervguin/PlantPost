FROM golang:1.24-alpine
LABEL AUTHOR="Impervguin"

RUN apk add --no-cache make curl libstdc++ libgcc

RUN mkdir /logs
RUN mkdir /build
WORKDIR /build

COPY go.* .
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go get -tool github.com/a-h/templ/cmd/templ@latest

RUN curl -LO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64-musl
RUN mv tailwindcss-linux-x64-musl /usr/local/bin/tailwindcss
RUN chmod +x /usr/local/bin/tailwindcss

COPY ./cmd/api/ ./cmd/api/
COPY ./internal/ ./internal/
COPY ./config/*.yaml ./config/
COPY tailwind.config.js .
COPY ./makefile .

CMD ["make", "api"]

