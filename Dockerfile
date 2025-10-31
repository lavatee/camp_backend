FROM golang:1.25.3-alpine
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go
CMD [ "./main" ]
EXPOSE 8000:8000