.PHONY: build clean deploy


build:
	cd api && env GOOS=linux go build -ldflags="-s -w" -o bin/create lib/create/main.go

clean:
	rm -rf ./api/bin

deploy: clean build
	sls deploy --verbose

goose:
	cd db && goose mysql "pingpen_user:$(PINGPENPW)@tcp(pingpen-dev-pingpendb-1my6lbs6pgl9t.cluster-cqwg3ufim27v.us-west-2.rds.amazonaws.com:3306)/pingpen?parseTime=true" $(filter-out $@,$(MAKECMDGOALS))


%:      
	@:    


