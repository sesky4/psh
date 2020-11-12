psh:
	go build -o build/psh psh/cmd/psh

pscp:
	go build -o build/pscp psh/cmd/pscp

psh_install: psh
	mv build/psh /usr/local/bin/psh

pscp_install: pscp
	mv build/pscp /usr/local/bin/pscp

clean:
	rm -rf build
