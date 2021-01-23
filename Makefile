ENV_TEST_FILE := .env.test
ENV_TEST = $(shell cat $(ENV_TEST_FILE))

.PHONY: up
up:
	docker-compose -f docker-compose.local.yml up -d --build

.PHONY: down
down:
	docker-compose -f docker-compose.local.yml down

.PHONY: logs
logs:
	docker-compose -f docker-compose.local.yml logs --tail=20

.PHONY: up-test-db
up-test-db:
	docker run --rm --env-file=$(ENV_TEST_FILE) -v $(PWD)/build/db/my.cnf:/etc/mysql/conf.d/my.cnf  --name blog-server_test_db_1 -d -p 3306:3306 mysql:8.0

.PHONY: down-test-db
down-test-db:
	docker stop blog-server_test_db_1
.PHONY: test

test:
	$(ENV_TEST) go test -v ./... -count=1


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