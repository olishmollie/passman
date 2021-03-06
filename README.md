# passman
A UNIX password manager

#### Installation
Passman is distributed via Homebrew.
```
$ brew install olishmollie/tools/passman
```

#### Getting started
Passman stores encrypted passwords in a directory in your home folder. To get started, run `passman init`, which will create the password store at `~/.passman` and generate an encryption key. 

Encrypt and store a password:
```
$ passman add foo/bar baz
$ passman add baz/boom bang
```
See your password store:
```
$ passman
Password Store
├── baz
│   └── boom
└── foo
    └── bar
```
Print a password to stdout:
```
$ passman foo/bar
baz
```
Or copy it to clipboard:
```
$ passman -c foo/bar
```
Dump encrypted passwords into backup file:
```
$ passman dump -o pswds.bak
```
Import a dump file:
```
$ passman import pswds.bak
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
	passman dump [-o <outfile>]
	passman edit <prefix>
	passman generate [-cn] [-l int]
	passman import <infile>
	passman init
	passman nuke [-f]
	passman -h | --help
	passman -v | --version

Options:
	-c, --copy                Copy to clipboard. 	  
	-f, --force               Nuke w/o confirmation.
	-h, --help                Show this screen.
	-l, --length=<int>        Specify length of generated password.
	-n, --nosym               Generate password w/ no symbols.
	-o, --out=<outfile>       Specify file to be written to [default: pswds~].
	-v, --version             Show version.
```