PROJECT_PATH=$(CURDIR)

lint: 
	docker run --rm --volume="${PROJECT_PATH}:/go/src/coinconv" -w /go/src/coinconv golangci/golangci-lint:v1.42-alpine golangci-lint run -E gofmt --skip-dirs=./vendor --deadline=10m

build:
	docker build . -t coinconv:local

.PHONY: $(command) $(sandbox) $(api_key)
run-docker: 
	docker run --env COINCONV_SANDBOX=$(sandbox) --env COINCONV_API_KEY=$(api_key) coinconv:local $(command)

