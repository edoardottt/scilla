REPO=github.com/edoardottt/scilla

fmt:
	@gofmt -s ./*; \
	echo "Done."

remod:
	rm -rf go.*
	go mod init ${REPO}
	go get
	echo "Done."

update:
	@go get -u; \
	go mod tidy -v; \
	echo "Done."

linux:
	@go build -o ./scilla
	mv ./scilla /usr/bin/
	cp -r lists/ /usr/bin/
	echo "Done."

unlinux:
	rm -rf /usr/bin/scilla
	rm -rf /usr/bin/lists/
	echo "Done."

test:
	@go test -v -race ./... ; \
	echo "Done."
