# set default shell
SHELL = bash -e -o pipefail

default: build

help:
	@echo "Usage: make [<target>]"
	@echo "where available targets are:"
	@echo
	@echo "build             : Build ethertool binary"
	@echo "deploy            : deploy binary to VM"
	@echo "build-and-deploy  : Build and deploy"
	@echo "help              : Print this help"
	@echo

build:
	@echo "building ethertool into app/"
	mkdir -p app
	GOOS=linux GOARCH=amd64 go build -o app/ethertool cmd/ethertool/main.go

deploy:
	@echo "deploy tool to my VM"
	scp -i id_rsa app/ethertool maik@192.168.56.101:/home/maik/test-tool/.

build-and-deploy: build deploy

