bin/ec2:
	cd ec2 && GOOS=linux CGO_ENABLED=0 go build -o ../bin/ec2

clean:
	rm bin/*

all: bin/ec2
