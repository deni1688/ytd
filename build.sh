export GO111MODULE=on
export GOOS=linux
export GOARCH=amd64

go build -tags netgo -o ytd

docker build -t deni1688/ytd .

docker push deni1688/ytd
