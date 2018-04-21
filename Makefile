AWS_ACCOUNT := $(shell aws sts get-caller-identity | jq -r .Account)
AWS_REGION := $(shell aws configure get region)

CMD := key-pair

STACK_NAME := Custom-KeyPair
STACK_BUCKET := cfn-$(AWS_ACCOUNT)-$(AWS_REGION)
STACK_PREFIX := resources/ec2/$(CMD)

GO := GOOS=linux GOARCH=amd64 go

build: ./target/$(CMD) Gopkg.lock
	
test:
	go test -v  ./...

clean:
	@rm -rf target packaged.yaml

bucket:
	aws s3 mb s3://$(STACK_BUCKET)

package: target/template.yaml

deploy:
	aws cloudformation deploy \
		--capabilities CAPABILITY_IAM \
		--stack-name $(STACK_NAME) \
		--template-file target/template.yaml

.PHONY: bucket build clean deploy package test

target/$(CMD): $(shell find aws cmd -name *.go)
	$(GO) build -v -o $@ ./cmd/$(CMD)

target/template.yaml: template.yaml build
	aws cloudformation package \
		--template-file $< \
		--output-template-file $@ \
		--s3-bucket $(STACK_BUCKET) \
		--s3-prefix $(STACK_PREFIX)

Gopkg.lock: Gopkg.toml
	dep ensure -v
