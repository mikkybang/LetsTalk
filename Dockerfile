FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY views views
COPY values values
COPY model model
COPY controller controller
COPY config.json .
COPY main.go .
RUN go build
ENTRYPOINT ./LetsTalk
EXPOSE 8080