FROM       golang:1.15.2
WORKDIR    /worker
COPY       go.mod .
COPY       go.sum .
RUN        go mod download
COPY       . .
RUN        go build -o worker ./cmd/main.go
ENTRYPOINT ["./registry"]
