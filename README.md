# PSH - autofill ssh password for you

## feature
 - [x] auto-fill password
 - [x] glob connection name match
 - [ ] other ssh options
 - [x] scp support

## How to use

1.write your own pconfig file in ~/.ssh/pconfig(example ./pconfig in repo)

2.make psh && make install_psh

3.make pscp && make install_pscp

4.(first time) `COMP_INSTALL=yes psh` to install bash/zsh auto_completion

5.Use it!
```
psh server1

scp server1:/root/my_file.py .
scp ./my_file.py server1:/root/
```
