deploy:
	cd server; \
	GOOS=linux go build ./main.go; \
	zip main.zip ./main; \
	aws lambda update-function-code \
	  --function-name pingpen_add \
	  --zip-file fileb://main.zip; \
	aws lambda update-function-code \
	  --function-name pingpen_sub \
	  --zip-file fileb://main.zip 
new_lambda:
	cd server; \
	GOOS=linux go build ./main.go; \
	zip main.zip ./main; \
	aws lambda create-function \
	  --region us-west-2 \
	  --function-name pingpen_$(FUNC) \
	  --memory 128 \
	  --role arn:aws:iam::532898105683:role/LambdaDeployer \
	  --runtime go1.x \
	  --zip-file fileb://main.zip \
	  --handler main 
	  #--environment Variables="{function=$(FUNC)}"
pb: server/pb/service.proto
	protoc -I$(GOPATH)/src/github.com/google/protobuf/src -I$(GOPATH)/src/github.com/mikerjacobi/pingpen/server/pb --go_out=./server/pb $(GOPATH)/src/github.com/mikerjacobi/pingpen/server/pb/*.proto
mysql:
	mysql -hnotify.cs9ds6yfnikc.us-east-1.rds.amazonaws.com -udbuser -p$(shell cat /etc/secrets/notify-db.json | grep password | cut -d'"' -f4) -Dnotify

