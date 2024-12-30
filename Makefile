compose:
	@docker-compose up up --build -f compose.dev.yaml

compose-down:
	@docker-compose up down -f compose.dev.yaml

container:
	@docker build -t subrotokumar/lbx . -f prod.Dockerfile