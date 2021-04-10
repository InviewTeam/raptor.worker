FROM       golang:1.15.2
# Configure libs for ffmpeg
RUN        apt-get update
RUN        apt-get install -y build-essential cmake libgtk2.0-dev pkg-config libavcodec-dev libavformat-dev libswscale-dev  libavresample-dev ffmpeg
WORKDIR    /worker
COPY       go.mod .
COPY       go.sum .
RUN        go mod download
COPY       . .
RUN        go build -o worker ./cmd/main.go
ENTRYPOINT ["./worker"]


