# passman
A UNIX password manager

#### Installation
```
$ go get github.com/olishmollie/passman
```

#### Getting started
Passman stores encrypted passwords in a directory in your home folder. To get started, run `passman init`, which will create the password store at `~/.passman` and generate an encryption key. 

Store a password like this:
```
passman add Email/myaccount secret
```
And copy it to your clipboard like this:
```
passman -c Email/myaccount
```

#### Tab Completion
Passman comes bundled with a bash completion script. To use it, source passman-completion.bash from your bashrc or bash_profile.

#### Usage
```
Usage:
	passman
	passman [-c] <prefix>
	passman add <prefix> <password>
	passman delete <prefix>
	passman dump
	passman edit <prefix>
	passman generate [-cn] [-l int]
	passman import <infile>
	passman init
	passman lock
	passman unlock
	passman -h | --help
	passman -v | --version

Options:
	-h, --help               Show this screen.
	-v, --version            Show version.
	-c, --copy               Copy to clipboard. 	  
	-n, --nosym              Generate password w/ no symbols.
	-l int, --length=int     Specify length of generated password.
```