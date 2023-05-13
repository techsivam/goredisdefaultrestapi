build:
	docker-compose build

run:
	docker-compose up -d

stop:
	docker-compose down

logs:
	docker-compose logs -f

test:
	curl -i http://localhost:8080
	echo
	curl -i http://localhost:8080/ping
	echo
	curl -X POST -H "Content-Type: application/json" -d '{"key":"sample","value":"example"}' http://localhost:8080/tenant1
	echo
	curl -i http://localhost:8080/tenant/tenant1
	echo

.PHONY: build run stop logs test
