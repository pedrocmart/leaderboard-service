.PHONY: test
test:
	@go test -race -cover ./...

.PHONY: local
local:
	make dev
	make docker-build
	make composeup

.PHONY: dev
dev:
	CGO_ENABLED=0 go build -o dist/leaderboard-service main.go

.PHONY: devlinux
devlinux:
	@GOOS=linux CGO_ENABLED=0 go build -o dist/leaderboard-service main.go

.PHONY: docker-build
docker-build:
	@docker build -t leaderboardservice .

.PHONY: composeup
composeup:
	docker-compose -f docker-compose.yml up && make composedown

.PHONY: composedown
composedown:
	docker-compose -f docker-compose.yml rm -s
