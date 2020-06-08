psh:
	go build cmd/psh/* -o build/psh

pscp:
	go build cmd/pscp/* -o build/pscp

install:
	mv build/psh /usr/local/bin/psh
	mv build/pscp /usr/local/bin/pscp

clean:
	rm -rf build
