REPO=github.com/edoardottt/scilla

fmt:
	@gofmt -s ./*
	@echo "Done."

remod:
	@rm -rf go.*
	@go mod init ${REPO}
	@go get
	@echo "Done."

update:
	@go get -u
	@go mod tidy -v
	@echo "Done."

lint:
	@golangci-lint run

linux:
	@go build ./cmd/scilla
	@sudo mv scilla /usr/bin/
	@sudo cp -r lists/ /usr/bin/
	@chmod +x scripts/config.sh
	@./scripts/config.sh
	@echo "Done."

unlinux:
	@sudo rm -rf /usr/bin/scilla
	@sudo rm -rf /usr/bin/lists/
	@rm -rf ~/.config/scilla
	@echo "Done."

test:
	@go test -v -race ./...
	@echo "Done."
