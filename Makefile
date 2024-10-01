logsdb:
	docker run -d --name logsdb -p 27017:27017 \
	-e MONGO_INITDB_ROOT_USERNAME=user \
	-e MONGO_INITDB_ROOT_PASSWORD=password \
	-v mongo_data:/data/db \
	mongo

build:
	go build

start: build
	sudo ./filechangestracker

test: 
	go test -v -cover ./...