FROM golang:1.25-alpine AS base
WORKDIR /app
RUN apk --no-cache git
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
COPY . .
EXPOSE  3000
CMD [ "go","run","main.go" ]

FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux -o bin/main ./main.go

FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/bin/main .
EXPOSE 3000
RUN ./main