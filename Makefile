ccgreen=\033[92m
ccred=\033[0;31m
ccyellow=\033[0;33m
ccend=\033[0m

test:
	@echo "$(ccgreen)Running unit tests$(ccend)"
	go test -coverprofile=/tmp/profile.out ./...

test-report:
	@echo "$(ccgreen)Generating unit test report$(ccend)"
	go test -coverprofile=/tmp/profile.out ./... && go tool cover -html=/tmp/profile.out
