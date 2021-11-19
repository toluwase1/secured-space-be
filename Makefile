
start-db:
	brew services start mongodb-community@4.2

stop-db:
	brew services stop mongodb-community@4.2

mock-db:
	mockgen -source=db/db.go -destination=db/mock_db.go -package=db

test: mock-db
	GIN_MODE=test go test ./...

run: |
	gofmt -w .
	go run main.go 

build:
	docker build . -t rentals -f Dockerfile

up:
	docker-compose up --build
	
# docker run -it --name=rentals-api -p 8080:8080 rentals