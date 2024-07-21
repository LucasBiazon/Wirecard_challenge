FROM golang:1.22.4-alpine

WORKDIR /wirecard

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR /wirecard/cmd/wirecard

RUN go build -o /wirecard/wirecard

EXPOSE 8080

CMD ["/wirecard/wirecard"]