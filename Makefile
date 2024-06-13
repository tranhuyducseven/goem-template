.PHONY: test coverage clean run compose down watch

test:
	ENV=test go test -coverprofile=cover.out -v ./...
coverage:
	go tool cover -html=cover.out
clean:
	rm main cover.out || true
	docker compose down --volumes --remove-orphans
down:
	docker compose down --remove-orphans
run:
	go run cmd/api/main.go
watch:
	air -c .air.toml
compose:
	docker compose -f compose.yml up -d --remove-orphans
start:
	docker compose -f app.compose.yml up -d --remove-orphans
stop:
	docker compose -f app.compose.yml down --remove-orphans
restart:
	docker compose -f app.compose.yml down --remove-orphans
	docker compose -f app.compose.yml up -d --remove-orphans
shutdown:
	docker compose -f app.compose.yml down --remove-orphans
	docker compose down --remove-orphans