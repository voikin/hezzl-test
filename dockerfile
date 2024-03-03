FROM golang:latest as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux GOARH=amd64 go build -o api cmd/api/main.go


FROM golang

WORKDIR /bin/app

COPY --from=builder /app/api api

EXPOSE 8080
ENTRYPOINT ["/bin/app/api"]
