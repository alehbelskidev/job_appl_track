FROM golang:1.25-alpine
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum .air.toml ./
RUN go mod download
EXPOSE 3001
CMD ["air"]
