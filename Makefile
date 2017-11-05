build:
	go build -ldflags "-X main.build=`date -u +%Y%m%d.%H%M%S`" -o svxm