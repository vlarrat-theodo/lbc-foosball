.PHONY: clean
clean:
	rm -rf __binaries/*

.PHONY: install
install:
	go mod tidy
	go install github.com/gobuffalo/pop/soda
	$(MAKE) launch-database
	sleep 10    # wait for DB to be up
	$(MAKE) migrate

.PHONY: launch-database
launch-database:
	docker-compose up -d

.PHONY: build
build:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o __binaries/scores/StoreGoal ./app/scores/StoreGoal

.PHONY: check_upx
check_upx:
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
compress:
	$(MAKE) check_upx
	upx --brute __binaries/scores/StoreGoal

.PHONY: local-deploy
local-deploy:
	sam local start-api --env-vars env.json --docker-network host

.PHONY: start
start: clean build launch-database local-deploy

.PHONY: aws-package
aws-package:
	sam package --output-template-file packaged.yaml --s3-bucket lbc-foosball-code --profile lbc-foosball

.PHONY: aws-deploy
aws-deploy:
	sam deploy --template-file packaged.yaml --stack-name lbc-foosball --capabilities CAPABILITY_IAM --profile lbc-foosball
	aws cloudformation describe-stacks --stack-name lbc-foosball --query 'Stacks[].Outputs' --profile lbc-foosball

.PHONY: deploy
deploy: clean build compress aws-package aws-deploy

.PHONY: test
test:
	go test ./...

.PHONY: migrate
migrate:
	soda migrate up
