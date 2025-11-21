DOCKER_COMPOSE_BIN = docker compose

run-release:
	source populate_variables && \
	${DOCKER_COMPOSE_BIN} -f deploy/release/docker-compose.yml --project-directory . up -d --force-recreate --build

run-dev:
	source populate_variables && \
	${DOCKER_COMPOSE_BIN} -f deploy/develop/docker-compose.yml --project-directory . up -d --force-recreate --build

stop-dev:
	source populate_variables && \
	${DOCKER_COMPOSE_BIN} -f deploy/develop/docker-compose.yml --project-directory . down

func-test: run-dev
	cd src/test/functional && \
	go test -v . -count=1

unit-test:
	cd src && \
	go test -count=1 -v ./cmd #TODO

prepare-stress: run-dev
	cd src/test/stress && \
	go test -v . -count=1

k6_local_%: 
	docker run --network internal --rm -i grafana/k6:latest run - <src/test/stress/k6/scripts/$*

stress-test: prepare-stress
	make k6_local_constant1000rps.js
