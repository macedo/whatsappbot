GITHUB_USERNAME ?= macedo
APPLICATION_NAME ?= whatsappbot

VERSION=$(shell cat VERSION)

build:
	docker build --tag ghcr.io/${GITHUB_USERNAME}/${APPLICATION_NAME}:${VERSION} .

push:
	docker push ghcr.io/${GITHUB_USERNAME}/${APPLICATION_NAME}:${VERSION}

release:
	docker pull ghcr.io/${GITHUB_USERNAME}/${APPLICATION_NAME}:${VERSION}
	docker tag ghcr.io/${GITHUB_USERNAME}/${APPLICATION_NAME}:${VERSION} ghcr.io/${GITHUB_USERNAME}/${APPLICATION_NAME}:latest
	docker push ghcr.io/${GITHUB_USERNAME}/${APPLICATION_NAME}:latest

version:
	@echo $(VERSION)
