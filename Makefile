compose:
	@docker-compose up --build -f ./dev/compose.dev.yaml

compose-down:
	@docker-compose down -f ./dev/compose.dev.yaml

container:
	@docker build -t subrotokumar/lbx . -f prod.Dockerfile