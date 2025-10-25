run-release:
	source populate_variables && \
	docker-compose -f deploy/release/docker-compose.yml --project-directory . up -d --force-recreate --build

run-dev:
	source populate_variables && \
	docker-compose -f deploy/develop/docker-compose.yml --project-directory . up -d --force-recreate --build

stop-dev:
	source populate_variables && \
	docker-compose -f deploy/develop/docker-compose.yml --project-directory . down

func-test: run-dev
	cd src/test/functional && \
	go test -v .

unit-test:
	cd src && \
	go test -count=1 -v ./cmd #TODO

prepare-stress: run-dev
	cd src/test/stress && \
	go test -v . -count=1

k6_local_%: 
	docker run --network internal --rm -i grafana/k6:latest run - <src/test/stress/k6/scripts/$*

stress-test: prepare-stress
	sleep 20 && \
	make k6_local_constant1000rps.js
