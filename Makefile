setup/osquery/mac:
	@./setup-osquery.sh


logsdb:
	docker stop logsdb; docker rm logsdb; true;
	docker run -d --name logsdb -p 27017:27017 \
	-e MONGO_INITDB_ROOT_USERNAME=user \
	-e MONGO_INITDB_ROOT_PASSWORD=password \
	-v mongo_data:/data/db \
	mongo

build:
	go build

start: build
	echo "> staring osquery!"
	sudo osqueryd --verbose --disable_events=false --disable_audit=false --disable_endpointsecurity=false --disable_endpointsecurity_fim=false --enable_file_events=true > test.log 2>&1 & 

	echo "> staring app!"
	sudo -S nohup ./filechangestracker > test.log 2>&1 & 

test: 
	go test -v -cover ./...

stop:
	sudo pkill osqueryd
	sudo -S pkill -f ./filechangestracker