.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Delete built binaries
	rm -rf __binaries/*

.PHONY: install
install: ## Install project and dependencies
	go mod tidy
	go install github.com/gobuffalo/pop/soda
	$(MAKE) launch-database
	sleep 10    # wait for DB to be up
	$(MAKE) migrate

.PHONY: launch-database
launch-database: ## Launch local database
	docker-compose up -d

.PHONY: build
build: ## Build binaries for each Lambda function
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o __binaries/scores/StoreGoal/StoreGoal ./app/scores/StoreGoal
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o __binaries/scores/FetchUserBalance/FetchUserBalance ./app/scores/FetchUserBalance

.PHONY: check_upx
check_upx: ## Check if UPX is installed (used for binaries compression)
ifeq ($(shell which upx), )
	@echo 'UPX is missing: it will be installed on this system'
ifeq ($(shell uname -s), Linux)
	sudo apt-get install upx-ucl
else ifeq ($(shell uname -s), Darwin)
ifeq ($(shell which brew), )
	@echo 'brew is needed to install UPX but is missing: you must install it manually!'
	@echo '    https://brew.sh/'
	@exit 1
else
	brew install upx
endif
endif
else
	@echo 'UPX is already installed on this system'
endif

.PHONY: compress
compress: ## Compress binaries after building
	$(MAKE) check_upx
	upx --brute __binaries/scores/StoreGoal/StoreGoal
	upx --brute __binaries/scores/FetchUserBalance/FetchUserBalance

.PHONY: local-deploy
local-deploy: ## Launch Lambda functions locally
	sam local start-api --env-vars env.json --docker-network host

.PHONY: start
start: clean build launch-database local-deploy ## Start complete application locally

.PHONY: aws-package
aws-package: ## Send binaries to AWS S3 bucket
	sam package --output-template-file packaged.yaml --s3-bucket lbc-foosball-code --profile lbc-foosball

.PHONY: aws-deploy
aws-deploy: ## Create or update all needed resources in AWS through CloudFormation
	sam deploy --template-file packaged.yaml --stack-name lbc-foosball --capabilities CAPABILITY_IAM --profile lbc-foosball
	aws cloudformation describe-stacks --stack-name lbc-foosball --query 'Stacks[].Outputs' --profile lbc-foosball

.PHONY: deploy
deploy: clean build compress aws-package aws-deploy ## Deploy complete application in production

.PHONY: test
test: ## Launch tests with coverage
	go test ./... -cover

.PHONY: migrate
migrate: ## Launch DB migrations
	soda migrate up
