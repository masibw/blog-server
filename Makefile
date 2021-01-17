.PHONY: up
up:
	docker-compose -f docker-compose.local.yml up -d --build

.PHONY: down
down:
	docker-compose -f docker-compose.local.yml down

.PHONY: logs
logs:
	docker-compose -f docker-compose.local.yml logs --tail=20

.PHONY: test
test:
	go test -v ./...

.PHONY: deploy
deploy: prod_down prod_update prod_up

.PHONY: prod_update
prod_update:
	git pull origin main

.PHONY: prod_up
prod_up:
	docker-compose -f docker-compose.prod.yml up -d --build

.PHONY: prod_down
prod_down:
	docker-compose -f docker-compose.prod.yml down