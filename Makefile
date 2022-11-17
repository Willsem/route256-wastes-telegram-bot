CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/cmd/bot

all: format build test lint

build: bindir
	@go build -o ${BINDIR}/bot ${PACKAGE}

test:
	@go test ./...

run:
	@go run ${PACKAGE} -config config.example.yaml | pino-pretty

generate:
	@go generate ./...

lint: install-lint
	@${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	@mkdir -p ${BINDIR}

format: install-smartimports
	@${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	@test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	@test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

docker-run:
	@sudo docker compose up
