.PHONY: start_server
start_server:
	docker-compose up

.PHONY: start_server
remove_server:
	docker-compose stop
	docker-compose rm -f

.PHONY: format
format:
	go fmt ./...
	gci -w .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: unit_tests
unit_tests:
	go test ./... -tags=unit,!integration

.PHONY: integration_tests
integration_tests:
	go test ./... -tags=integration,!unit -count=1

.PHONY: coverage
coverage:
	go test ./... -tags=integration,unit -cover -count=1

.PHONY: test
test: format lint coverage

