FROM golang:1.24-alpine
LABEL AUTHOR="Impervguin"

RUN apk add --no-cache make curl libstdc++ libgcc

# TYPESCRIPT
RUN apk add --no-cache nodejs npm
RUN npm install -g typescript esbuild

# DIRS
RUN mkdir /logs
WORKDIR /build

# HOT RELOAD
RUN go install github.com/air-verse/air@latest

# GO PACKAGES
COPY go.* .
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go get -tool github.com/a-h/templ/cmd/templ@latest

# TAILWIND CSS
RUN curl -LO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64-musl
RUN mv tailwindcss-linux-x64-musl /usr/local/bin/tailwindcss
RUN chmod +x /usr/local/bin/tailwindcss

CMD ["air"]

