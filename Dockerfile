FROM golang:1.23.1-alpine3.20

WORKDIR /github-secret-synchronizer

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /github-secret-synchronizer/github-secret-synchronizer

CMD ["/github-secret-synchronizer/github-secret-synchronizer"]
