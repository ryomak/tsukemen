.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fixtures && sudo docker-compose up --force-recreate -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fixtures && sudo docker-compose down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./tsukemen

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/heroes-service-* heroes-service
	@sudo docker rm -f -v `sudo docker ps -a --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@sudo docker rmi `sudo docker images --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"
