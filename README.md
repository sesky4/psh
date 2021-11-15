# PSH - Automate SSH login

## feature
 - [x] auto-fill password
 - [x] glob connection name match
 - [x] other ssh options
 - [x] scp support

## How to use

1.write your own pconfig file in ~/.ssh/pconfig, example: ./pconfig.example

2.make psh && make install_psh

3.make pscp && make install_pscp

4.(first time) `psh -complete && pscp -complete` to install bash/zsh auto_completion, this step will append a few lines to your ~/.bashrc or ~/.zshrc

5.Use it !
```
psh server1

pscp server1:/root/file_to_download.py .
pscp ./file_to_upload.py server1:/root/
```
