#!/bin/bash
processName='recommend-server'
env=$1
if [ ! -n "$env" ];then
  echo "need input environment value:[dev/test/prd]"
  exit 1
fi
export REC_ENV="$env"
git pull origin main
git config --global http.extraheader "PRIVATE-TOKEN:glpat-yysVywQqx27fpCTM6Zcc"
git config --global url."git@gitlab.com:cher8/lion.git".insteadOf "https://gitlab.com/cher8/lion.git"
#chmod 600 ~/.ssh/id_rsa.pub && eval "$(ssh-agent -s)"
#ssh-add ~/.ssh/id_rsa.pub
export GOPROXY=https://goproxy.cn
export GOPRIVATE=gitlab.com/cher8
echo "GOPROXY="$GOPROXY
echo "GOPRIVATE="$GOPRIVATE
echo "generate protobuf code..."
protoc -I. --go_out=. --go-grpc_out=. apis/recommend.proto
if [ $? -eq 0 ];then
  echo "generate pd/grpc is success by recommend.proto"
else
  echo "generate pd/grpc is fail by recommend.proto"
  exit 1
fi
protoc -I. --go_out=. --go-grpc_out=. internal/sort/estimate/estimate.proto
if [ $? -eq 0  ];then
  echo "generate pd/grpc is success by estimate.proto"
else
  echo "generate pd/grpc is fail by estimate.proto"
  exit 1
fi
echo "Start compile..."
go mod tidy && go mod vendor && go build -o recommend-server cmd/main.go
if [ $? -eq 0 ]
then
    echo "go compile success"
else
    echo "go compile fail"
    exit 1
fi
PID=$(ps -ef|grep $processName|grep -v grep|awk '{printf $2}')
if [ $? -eq 0 ] && [ ${#PID} -gt 0 ]; then
    echo "process id : $PID"
    kill  ${PID}
    sleep 11
    if [ $? -eq 0 ]; then
        echo "kill $processName success"
    else
        echo "kill $processName fail"
        exit 1
    fi
else
    echo "process $processName not exist"
fi
./recommend-server -conf="cmd/$env".yaml &