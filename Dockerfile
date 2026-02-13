FROM golang:1.24-alpine

RUN apk add --no-cache git tzdata

WORKDIR /app

RUN go install github.com/air-verse/air@v1.52.3

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]