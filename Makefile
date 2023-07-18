


ifndef STACK_NAME
override STACK_NAME = url-proxy-test
endif

TEST_REGION = eu-west-1
TEST_ENV = test

NAME = $(STACK_NAME)

TEST_PARAMS = Name=$(NAME)

STACK_DIR = ./stack

stack/build:
	@cd $(STACK_DIR); \
	sam build --cached

stack/validate:
	@cd $(STACK_DIR); \
	sam validate --lint 

stack/package:
	@cd $(STACK_DIR); \
	sam package --output-template-file packaged.yaml --s3-bucket ln80-sam-pkgs

integ/deploy:
	@cd $(STACK_DIR); \
	sam deploy \
		--no-confirm-changeset \
		--no-fail-on-empty-changeset \
		--stack-name $(STACK_NAME) \
		--config-env $(TEST_ENV) \
		--capabilities CAPABILITY_IAM\
		--region $(TEST_REGION) \
		--parameter-overrides $(TEST_PARAMS) \
		--template-file packaged.yaml

