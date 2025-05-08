FROM golang:1.23

WORKDIR $GOPATH/src/github.com/edoardottt/scilla

COPY . .
RUN go mod download golang.org/x/sys
RUN go install -v ./...

ENTRYPOINT ["scilla"]