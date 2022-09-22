FROM golang:alpine
WORKDIR $GOPATH/src/test-go-liveness
ADD . ./
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
RUN go build -o test-go-liveness .
EXPOSE 8001
ENTRYPOINT  ["./test-go-liveness"]
