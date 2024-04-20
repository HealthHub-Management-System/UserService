FROM golang:alpine

WORKDIR /backend
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g ./cmd/api/main.go
RUN go build -o ./bin/api ./cmd/api \
    && go build -o ./bin/migrate ./cmd/migrate

CMD ["/backend/bin/api"]
EXPOSE 8080
