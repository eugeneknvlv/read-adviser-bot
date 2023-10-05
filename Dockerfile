FROM golang:latest
LABEL authors="Konovalov"

WORKDIR /read-adviser-bot

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o read-adviser-bot

ENTRYPOINT ["./read-adviser-bot"]