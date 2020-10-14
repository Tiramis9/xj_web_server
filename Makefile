#/bin/bash
# This is how we want to name the binary output
# 网关模块
OUTPUT=bin/xj_web_server
SRC=cmd/main.go
# account模块
SRCPROTO=module/account/main.go
# db 模块
SRCPROTO2=module/db/main.go


# public 模块
SRCPROTO3=module/public/main.go

# account模块 输出
PROTOPROTO=module/bin/server_account
PROTOPROTO2=module/bin/server_db
PROTOPROTO3=module/bin/server_public

# These are the values we want to pass for Version and BuildTime
GITTAG=1.0.0
BUILD_TIME=`date +%Y%m%d%H%M%S`



# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=${GITTAG} -X main.Build_Time=${BUILD_TIME} -s -w"

local:
	rm -f ./bin/xj_web_server_mac
	go build ${LDFLAGS} -o ${OUTPUT}_mac ${SRC}

server:
	rm -f ./bin/xj_web_serverr
	rm -f ./module/bin/server_*
	go build ${LDFLAGS} -o ${OUTPUT} ${SRC}
	go build ${LDFLAGS} -o ${PROTOPROTO} ${SRCPROTO}
	go build ${LDFLAGS} -o ${PROTOPROTO2} ${SRCPROTO2}
	go build ${LDFLAGS} -o ${PROTOPROTO3} ${SRCPROTO3}
proto:
	protoc --proto_path=${GOPATH}/src:. --go_out=. httpserver/wss/proto/*.proto
	##protoc --proto_path=module/account/proto   --go_out=module/account/proto --micro_out=module/account/proto module/account/proto/user.proto
	## protoc --proto_path=module/db/proto   --go_out=module/db/proto --micro_out=module/db/proto module/db/proto/db.proto
	## protoc --proto_path=module/public/proto   --go_out=module/public/proto --micro_out=module/public/proto module/public/proto/public.proto
#./protoc.sh

web-server:
	rm -f ./bin/xj_web_serverr
	rm -f ./module/bin/server_*
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter ${LDFLAGS} -o ${OUTPUT} ${SRC}
	# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter ${LDFLAGS} -o ${PROTOPROTO} ${SRCPROTO}
	# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter ${LDFLAGS} -o ${PROTOPROTO2} ${SRCPROTO2}
	# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter ${LDFLAGS} -o ${PROTOPROTO3} ${SRCPROTO3}


clean:
	rm -f ./bin/xj_web_server_*
	rm -f ./module/bin/server_*
