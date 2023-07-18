


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

integ-test/deploy:
	@cd $(STACK_DIR); \
	sam deploy \
		--no-confirm-changeset \
		--no-fail-on-empty-changeset \
		--stack-name $(STACK_NAME) \
		--config-env $(TEST_ENV) \
		--capabilities CAPABILITY_IAM\
		--region $(TEST_REGION) \
		--parameter-overrides $(TEST_PARAMS) \
		--template-file template.yaml

integ-test/run:
	@cd $(STACK_DIR); \
	\
	proxyUrl="`aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME) \
		--region $(TEST_REGION) \
		--query "Stacks[0].Outputs[?OutputKey=='ProxyFunctionUrl'].OutputValue" --output text`"; \
	\
	if [ $$? != 0 ]; then echo "fetch cfn resource failed"; exit 2; fi; \
	if [ -z "$$proxyUrl" ]; then echo "invalid cfn output value"; exit 2; fi; \
	\
	PROXY_FUNCTION_URL=$$proxyUrl go test --tags=integ -race -cover -v
	
integ-test/clear:
	@cd $(STACK_DIR); \
	sam delete --no-prompts --stack-name $(STACK_NAME) --region $(TEST_REGION)

