CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin

LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}

PACKAGE=gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/cmd/bot

PACKAGE_MIGRATE=gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/cmd/migrate
MIGRATIONS_DIR=file://internal/migrations
DATABASE_URL=postgresql://${DATABASE_USER}:${DATABASE_PASS}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_DB}?sslmode=disable

all: build test lint

build: bindir
	@go build -o ${BINDIR}/bot ${PACKAGE}

bindir:
	@mkdir -p ${BINDIR}

test:
	@go test ./...

run:
	@go run ${PACKAGE} -config configs/config.example.yaml | pino-pretty

generate:
	@go generate ./...
	@go run ${PACKAGE_MIGRATE} -config configs/config.migrate.yaml -name ${MIGRATION_NAME}

lint: install-lint
	@${LINTBIN} run

install-lint: bindir
	@test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

precommit: build test lint
	echo "OK"

docker-run:
	@sudo docker compose up -d

migrate:
	@atlas migrate apply --dir "${MIGRATIONS_DIR}" --url "${DATABASE_URL}"