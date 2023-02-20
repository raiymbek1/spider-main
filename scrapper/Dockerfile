FROM golang:alpine AS scrapper
WORKDIR /app
RUN apk add --no-cache git

COPY . .

RUN go build -o app ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=scrapper /app/app .

COPY conf.env .

EXPOSE 8080

CMD ["./app"]