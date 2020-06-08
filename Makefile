psh:
	go build -o build/psh cmd/psh/*

pscp:
	go build -o build/pscp cmd/pscp/*

install:
	mv build/psh /usr/local/bin/psh
	mv build/pscp /usr/local/bin/pscp

clean:
	rm -rf build
