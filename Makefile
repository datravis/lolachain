GOARCH = amd64

PACKAGE_DIR=$github.com/datravis/lolachain/cmd

# Build the project
all:  deps lint test linux darwin windows

linux: 
	GOOS=linux GOARCH=${GOARCH} go build -o dist/linux/lolachain-validator github.com/datravis/lolachain/cmd/validator
	GOOS=linux GOARCH=${GOARCH} go build -o dist/linux/lolachain-wallet github.com/datravis/lolachain/cmd/wallet
	GOOS=linux GOARCH=${GOARCH} go build -o dist/linux/lolachain-gui github.com/datravis/lolachain/cmd/gui
darwin:
	GOOS=darwin GOARCH=${GOARCH} go build -o dist/darwin/lolachain-validator github.com/datravis/lolachain/cmd/validator
	GOOS=darwin GOARCH=${GOARCH} go build -o dist/darwin/lolachain-wallet github.com/datravis/lolachain/cmd/wallet
	GOOS=darwin GOARCH=${GOARCH} go build -o dist/darwin/lolachain-gui github.com/datravis/lolachain/cmd/gui

windows:
	GOOS=windows GOARCH=${GOARCH} go build -o dist/windows/lolachain-validator github.com/datravis/lolachain/cmd/validator
	GOOS=windows GOARCH=${GOARCH} go build -o dist/windows/lolachain-wallet github.com/datravis/lolachain/cmd/wallet
	GOOS=windows GOARCH=${GOARCH} go build -o dist/windows/lolachain-gui github.com/datravis/lolachain/cmd/gui

lint:
	-gometalinter ./... --exclude=vendor

test:
	go test ./...

deps:
	dep ensure --vendor-only


.PHONY: linux darwin windows test deps