FROM golang:1.23.4-alpine3.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gss .

FROM alpine:3.21.0

WORKDIR /app

COPY --from=build /app/gss .

RUN apk add --no-cache ca-certificates tzdata

CMD ["/app/gss"]
