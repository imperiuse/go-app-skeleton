.EXPORT_ALL_VARIABLES:

-include .env

APP_NAME ?= go-app-skeleton
APP_ENV ?= dev
APP_VERSION ?= dev

# $CI_REGISTRY ?= registry.gitlab.bla-bla-bla
# GOLANG_IMAGE ?= ${CI_REGISTRY}/golang:latest
# ALPINE_IMAGE ?= ${CI_REGISTRY}/docker-images/alpine:latest
# NEXUS_HOST ?= nexus.bla-bla-bla
# NEXUS_PORT ?= 8081
# NEXUS_USER ?= some_user
# NEXUS_PSWD ?= <PASSWORD>

#BUILD_WITH_DEBUG    ?= CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${NAME} -gcflags "all=-N -l" -ldflags '-v       -linkmode internal -extldflags \"-static\" -X ${GO_PACKAGE}/app/config.Version=${VERSION}' ${GO_PACKAGE}
#BUILD_WITHOUT_DEBUG ?= CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${NAME}
# CGO_ENABLE=0 -  важен!!!! бинарь без него был кропнутый -> https://issue.life/questions/34729748
# We have to add an environment variable CGO_ENABLED=0 to disable dynamically links for a few dependencies. Normally we could not run Go applications from scratch because of this. We can also get rid of two more things in our binary. DWARF tables and annotations. The tables are needed for debuggers and the annotations for stack traces. Adding-ldflags="-s -w" removes them from our binary
BUILD_CMD ?= CGO_ENABLED=0 go build -tags=jsoniter -a -v -o bin/${APP_NAME} -ldflags '-v -w -s -linkmode auto -extldflags \"-static\" -X  main.AppName=${APP_NAME}  -X  main.AppVersion=${APP_VERSION}  -X  main.AppEnv=${APP_ENV}' ./cmd/${APP_NAME}
UPX_CMD ?= upx --best --lzma bin/${APP_NAME}

MACHINE_IP ?= 127.0.0.1
.EXPORT_ALL_VARIABLES:
POSTGRESQL_URL = 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_USER}?sslmode=disable'

.PHONY: git_submodule_update
git_submodule_update:
	git submodule init
	git submodule update
	git submodule foreach git pull origin master

.PHONY: go_mod_upgrade
go_mod_upgrade:
	cd ./cmd/${APP_NAME}/ && go get -u && cd -

.PHONY: docker_clean_all
docker_clean_all:
	yes | docker rm -f `docker ps -aq`
	yes | docker volume rm -f `docker volume ls -q`
	yes | docker network prune
	yes | docker system prune

.PHONY: dev_env_up
dev_env_up:
	@echo "1. Run docker-compose dev"
	docker-compose -f .docker/docker-compose.override.yml up -d
	@echo "2. Install pg extension"
	./tools/check_postgres_ready.sh localhost 5432
	psql postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_USER} -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
	@echo "3. Running migrations"
	docker run --rm  -v ${PWD}/migrations/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database ${POSTGRESQL_URL} up

.PHONY: dev_env_down
dev_env_down:
	@echo "Stop dev environment"
	docker-compose -f .docker/docker-compose.override.yml down

.PHONY: test_env_up
test_env_up:
	@echo "Start test environment"
	@echo "1. Export test env variables and Run docker-compose test"
	docker-compose -f .docker/docker-compose.test.yml up -d
	@echo "2. Install pg extension"
	./tools/check_postgres_ready.sh localhost 5433
	psql postgresql://${CORE_TEST_POSTGRES_USER}:${CORE_TEST_POSTGRES_PASSWORD}@localhost:5433/${CORE_TEST_POSTGRES_USER} -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
	@echo "3. Running migration"
	docker run --rm  -v ${PWD}/migrations/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database ${TEST_POSTGRESQL_URL} up

.PHONY: test_env_down
test_env_down:
	@echo "Stop test environment"
	docker-compose -f .docker/docker-compose.test.yml down

.PHONY: new_migration
new_migration:
	@echo "Create new migration files"
	touch ./migrations/migrations/`date +%s`_$(name).down.sql
	touch ./migrations/migrations/`date +%s`_$(name).up.sql

.PHONY: do_migrate
do_migrate:
	@echo "Running migration"
	docker run --rm -v ${PWD}/migrations/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database ${POSTGRESQL_URL} up

.PHONY: undo_migrate
undo_migrate:
	@echo "Running undo migration"
	echo y | docker run --rm -v ${PWD}/migrations/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database ${POSTGRESQL_URL} down 1

.PHONY: lint
lint:
	@echo "Running golangci-lint"
	docker run --rm -v /tmp/${APP_NAME}/.cache:/root/.cache -v /tmp/${APP_NAME}/linter/pkg:/go/pkg/mod/ -v ${PWD}:/${APP_NAME} -w=/${APP_NAME} registry-tc.dev.codd.local/theseus/golangci-lint:latest golangci-lint run -v --deadline 5m ./...

.PHONY: swagger_docs
swagger_docs:
	@echo "Create Swagger docs"
	swag init --parseDependency --parseInternal --parseDepth 1 -g cmd/${APP_NAME}/main.go
	@echo "You can see docs here -> http://localhost:8081/swagger/index.html"

.PHONY: tests
tests:
	@echo "Running go tests"
	go test -timeout 5m -v -race -short `go list ./... ` && cd -

.PHONY: coverage
coverage:
	@echo "Running coverage.sh script"
	./tools/coverage.sh

.PHONY: integration_tests
integration_tests:
	@echo "Running integrations tests"
	@echo "1. Run test env"
	make test_env_up
	sleep 1  # sleep driven development
	@echo "2. Run tests"
	go test -timeout 5m -race -short `go list ./... | grep postgres`  # only postgres dir
	@echo "3. Stop test env"
	make test_env_down

.PHONY: integration_tests_new_way
integration_tests_new_way:
	@echo "Running integration tests"
	export `grep -v '^#' .env | xargs` && cd app/tests/integration && TESTCONTAINERS_RYUK_DISABLED=true go test -race -v -timeout 15m -p 1 `go list ./...` && cd -

.PHONY: build
build:
	@echo "Running build"
	${BUILD_CMD}
	@echo "Running UPX (zip) binary"
	${UPX_CMD}

.PHONY: run_generators
run_generators:
	@echo "Running generator: models_method_generator.go"
	go run ./tools/generators/models_method_generator.go -s User,Role,UsersRole,Session,Program,RoadController,CalendarRule,Route,AuditErrorsRCxx,AuditStateRCxx,AuditCommandsRCxx -in ${PWD}/internal/services/storage/model/model.go -out ${PWD}/internal/services/storage/model/model_gen.go

.PHONY: store_artifact_to_nexus
store_artifact_to_nexus:
	@echo "Store binary to Nexus"
	curl -v -u ${NEXUS_USER}:${NEXUS_PSWD} --upload-file filename ${NEXUS_HOST}-${APP_NAME}/${APP_NAME}_${APP_VERSION}

.PHONY: build_docker_local
build_docker_local:
	# NB! GOOD IDEA FOR DEBUG ALL DOCKER STUFF -> docker run -it --entrypoint /bin/sh <IMAGE_NAME>
	@echo build_docker_local
	docker build                                  \
    	--pull                                    \
        --build-arg GOLANG_IMAGE=${GOLANG_IMAGE}  \
        --build-arg ALPINE_IMAGE=${ALPINE_IMAGE}  \
        --build-arg NEXUS_HOST=${NEXUS_HOST}      \
        --build-arg NEXUS_PORT=${NEXUS_PORT}      \
        --build-arg NEXUS_USER=${NEXUS_USER}      \
        --build-arg NEXUS_PSWD=${NEXUS_PSWD}      \
        --build-arg APP_NAME=${APP_NAME}          \
        --build-arg APP_VERSION=${APP_VERSION}    \
        --build-arg APP_ENV=stage                 \
        --target app                              \
        -t core_local                             \
        -f .docker/Dockerfile .
#        --no-cache                               \

.PHONY: run_docker_compose_stage_local
run_docker_compose_stage_local:
	docker-compose -f .docker/docker-compose.override.yml -f ./.docker/docker-compose.stage.yml up


.PHONY: lint-local
lint-local:
	@echo "Run golangci-lint"
	cd app && golangci-lint run -v ./... && cd -
