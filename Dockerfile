
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

COPY . $GOPATH/src/telbotnsn/

WORKDIR $GOPATH/src/telbotnsn

RUN go mod tidy
# RUN go mod download
RUN mkdir /app
RUN go build -o /app/zbot .

FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/zbot /app/zbot
COPY --from=builder go/src/telbotnsn/.env /app/

WORKDIR /app 
CMD ["/app/zbot", "run"]