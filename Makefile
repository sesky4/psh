psh:
	go build

install: psh
	mv psh /usr/local/bin/psh

clean:
	rm psh