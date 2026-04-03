include .env
export

migrate-up: 
	migrate -path migrations -database ${CONN_STR} up

migrate-down:
	migrate -path migrations -database ${CONN_STR} down