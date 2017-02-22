all: clean deps build deploy

clean:
	rm -rf ./pipo

deps:
	glide install

build:
	export GOOS=linux && export GOARCH=arm && go build

deploy:
	scp ./pipo pi@$PIPO_TARGET
