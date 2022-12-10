FROM golang:1.20-rc-alpine3.17

WORKDIR /code
COPY . .

RUN apk add --no-cache gcc musl-dev linux-headers

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

RUN go mod vendor
RUN go mod download
RUN go build

EXPOSE 8080

CMD ["./application"]
