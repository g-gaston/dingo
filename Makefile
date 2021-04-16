.PHONY: build
build: image
	go build -o bin/dingo

image:
	docker pull public.ecr.aws/bacardi/alpine:3.13.0
	docker save public.ecr.aws/bacardi/alpine:3.13.0 -o alpine.tar