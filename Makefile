run-release:
	source populate_variables && \
	docker-compose -f deploy/release/docker-compose.yml --project-directory . up -d --force-recreate --build

run-dev:
	source populate_variables && \
	docker-compose -f deploy/develop/docker-compose.yml --project-directory . up -d --force-recreate --build

stop-dev:
	source populate_variables && \
	docker-compose -f deploy/develop/docker-compose.yml --project-directory . down

func-test:
	cd src/test/functional && \
	go test -v .

unit-test:
	cd src && \
	go test -count=1 -v ./cmd #TODO