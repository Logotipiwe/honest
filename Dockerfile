FROM golang:1.22-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag i -d ./src/cmd,./src/internal
RUN go build -v -o /usr/local/bin/app ./src/cmd

CMD ["app"]