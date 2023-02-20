FROM golang:alpine AS crawler
WORKDIR /app
RUN apk add --no-cache git

COPY . .

RUN go build -o app ./cmd/cli/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=crawler /app/app .

COPY conf.env .

CMD ["./app"]