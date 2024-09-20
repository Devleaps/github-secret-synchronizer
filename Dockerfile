FROM golang:1.23.1-alpine3.20

WORKDIR /github-secrets-synchronizer

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /github-secrets-synchronizer/github-secrets-synchronizer

CMD ["/github-secrets-synchronizer/github-secrets-synchronizer"]
