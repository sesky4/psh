psh:
	go build -o build/psh psh/cmd/psh

pscp:
	go build -o build/pscp psh/cmd/pscp

psh_install: psh
	go install ./cmd/psh

pscp_install: pscp
	go install ./cmd/pscp

install: psh_install pscp_install

complete:
	psh -complete -y ; pscp -complete -y

clean:
	rm -rf build
