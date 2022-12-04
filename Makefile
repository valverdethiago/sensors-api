infra_up:
	docker-compose -f ./docker/docker-compose.yml  up -d --build

infra_down:
	docker-compose -f ./docker/docker-compose.yml  down
