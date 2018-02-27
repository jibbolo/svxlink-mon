build:
	go build -ldflags "-X main.build=`date -u +%Y%m%d%H%M%S`-`git log --format=%h -1`" -o svxm
