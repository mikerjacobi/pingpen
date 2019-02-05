.PHONY: build clean deploy


build:
	cd api && env GOOS=linux go build -ldflags="-s -w" -o bin/create lib/create/main.go

clean:
	rm -rf ./api/bin

deploy: clean build
	sls deploy --verbose

sandbox: clean build
	./scripts/run_sandbox.py

devdb: 
	docker run --name pingpen_sandbox_db -e MYSQL_ROOT_PASSWORD=password -d mysql:5.6
	export DBHOST=$(docker inspect pingpen_sandbox_db  | grep 'IPAddress"' | head -n1 | cut -d":" -f2 | cut -d'"' -f2); \
	cat $(DBHOST) \ 
	cd db && goose mysql "root:password@tcp($(DBHOST):3306)/pingpen?parseTime=true" up
	
goose:
	cd db && goose mysql "pingpen_user:$(PINGPENPW)@tcp(pingpen-dev-pingpendb-1my6lbs6pgl9t.cluster-cqwg3ufim27v.us-west-2.rds.amazonaws.com:3306)/pingpen?parseTime=true" $(filter-out $@,$(MAKECMDGOALS))


%:      
	@:    


