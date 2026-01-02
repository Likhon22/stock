FROM golang:1.25-alpine AS base
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
COPY . .
CMD [ "go","run","main.go" ]

FROM base AS builder 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux -o go build bin/main ./main.go

FROM alpine:latest as prod
WORKDIR /app
COPY --from=builder /app/bin/main .
CMD [ "./main" ]