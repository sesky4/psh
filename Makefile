psh:
	go build

pscp:
	go build -o pscp pscp/*

install: psh pscp
	mv psh /usr/local/bin/psh
	mv pscp /usr/local/bin/pscp

clean:
	rm psh