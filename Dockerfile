FROM golang:1.16

WORKDIR $GOPATH/src/github.com/edoardottt/scilla

COPY . .
RUN go install -v ./...
RUN mv ./lists/ /usr/bin/

ENTRYPOINT ["scilla"]