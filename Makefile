TARGET=./build
ARCHS=amd64 386
LDFLAGS="-s -w"
GCFLAGS="all=-trimpath=$(shell pwd)"
ASMFLAGS="all=-trimpath=$(shell pwd)"

fmt:
	@gofmt -s ./*; \
	echo "Done."

remod:
	rm -rf go.*
	go mod init github.com/edoardottt/scilla
	go get

update:
	@go get -u; \
	go mod tidy -v; \
	echo "Done."

windows:
	@for GOARCH in ${ARCHS}; do \
		echo "Building for windows $${GOARCH} ..." ; \
		mkdir -p ${TARGET}/scilla-windows-$${GOARCH} ; \
		mv lists/ ${GOPATH}/bin; \
		GOOS=windows GOARCH=$${GOARCH} GO111MODULE=on CGO_ENABLED=0 go build -ldflags=${LDFLAGS} -gcflags=${GCFLAGS} -asmflags=${ASMFLAGS} -o ${TARGET}/scilla-windows-$${GOARCH}/scilla.exe ; \
	done; \
	echo "Done."

linux:
	@for GOARCH in ${ARCHS}; do \
		echo "Building for linux $${GOARCH} ..." ; \
		mkdir -p ${TARGET}/scilla-linux-$${GOARCH} ; \
		mv lists/ ${GOPATH}/bin; \
		GOOS=linux GOARCH=$${GOARCH} GO111MODULE=on CGO_ENABLED=0 go build -ldflags=${LDFLAGS} -gcflags=${GCFLAGS} -asmflags=${ASMFLAGS} -o ${TARGET}/scilla-linux-$${GOARCH}/scilla ; \
	done; \
	echo "Done."

test:
	@go test -v -race ./... ; \
	echo "Done."

clean:
	@rm -rf ${TARGET}/* ; \
	go clean ./... ; \
	echo "Done."