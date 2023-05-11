FROM golang:1.20-alpine as dev

ENV CGO_ENABLED=0
ENV RDB_DATABASE="chat"

RUN apk add --no-cache --update git tzdata ca-certificates
RUN apk add make

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest


COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["air", "-c", ".air.toml"]

# Build Image
FROM golang:1.20-alpine as build

ENV CGO_ENABLED=0
ENV RDB_DATABASE="chat"

RUN apk add --no-cache --update git tzdata ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Production Image (distroless) with only the binary
FROM gcr.io/distroless/static-debian11 as app
COPY --from=build /app/main /app/main
COPY --from=build /app/embedding /app/embedding

WORKDIR /app

CMD ["./main", "run"]