deploy:
	cd server
	GOOS=linux go build main.go -o lambda_handler
	zip handler.zip ./lamda_handler	
	aws lambda create-function \
	  --region us-west-2 \
	  --function-name lambda-handler \
	  --memory 128 \
	  --role arn:aws:iam::532898105683:role/LambdaDeployer \
	  --runtime go1.x \
	  --zip-file fileb://handler.zip \
	  --handler lambda-handler	
	cd ..
pb: server/pb/service.proto
	protoc -I$(GOPATH)/src/github.com/google/protobuf/src/google/protobuf -I$(GOPATH)/src/github.com/mikerjacobi/pingpen/server/pb --go_out=./server/pb $(GOPATH)/src/github.com/mikerjacobi/pingpen/server/pb/*.proto
mysql:
	mysql -hnotify.cs9ds6yfnikc.us-east-1.rds.amazonaws.com -udbuser -p$(shell cat /etc/secrets/notify-db.json | grep password | cut -d'"' -f4) -Dnotify

