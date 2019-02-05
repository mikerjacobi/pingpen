.PHONY: build clean deploy

#common commands

build:
	cd api && env GOOS=linux go build -ldflags="-s -w" -o bin/create lib/create/main.go

clean:
	rm -rf ./api/bin

# real env commands
deploy: clean build
	sls deploy --verbose

goose:
	cd db && goose mysql "pingpen_user:$(PINGPENPW)@tcp($(PINGPENDB):3306)/pingpen?parseTime=true" $(filter-out $@,$(MAKECMDGOALS))

#creates a fresh sandbox environment
sandbox: clean build
	docker stack deploy -c sandbox.yml pingpen

devgoose:
	cd db && goose mysql "root:password@tcp(127.0.0.1:3306)/pingpen?parseTime=true"  $(filter-out $@,$(MAKECMDGOALS))

devdeploy: clean build
	docker service scale -d pingpen_post_note=0
	docker service scale -d pingpen_post_note=1
%:      
	@:    


