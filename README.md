# passman
### A UNIX password manager

#### Installation
```
$ go get github.com/olishmollie/passman
$ go install github.com/olishmollie/passman
```

#### Getting started
Passman stores encrypted passwords in a directory in your home folder. To get started, run `passman init`, which will create the password store at `~/.passman` and ask you for a master passphrase, which will be used to generate an encryption key. This key is stored in `~/.passman/.fpubkey`.

#### Tab Completion
Passman comes bundled with a bash completion script. To use it, source passman-completion.bash from your bashrc or bash_profile.

#### Usage
```usage: passman [opts...] [command] [args...]
ex: passman touch Category/Website/username pswd
commands: dump edit generate import init rm touch
	passman - prints a tree of pswds in store

passman [opts...] <pswd_file> - prints unencrypted pswd
    opts:
        -copy - copies unencrypted pswd to clipboard

dump <outfile> - prints unencrypted pswds to outfile

edit <pswd_file> - edit pswd in editor set to $VISUAL

generate [opts...] - generates a random pswd
    opts:
        -copy - copies unencrypted pswd to clipboard
        -len=int - specifies length of generated pswd
        -nosym - generate a password with no symbols

import <infile> - imports passwords from infile.

init - create password store if it doesn't exist, and generate encryption key

rm <pswd_file> - remove <pswd_file> from pswd store

touch <pswd_file> - add <pswd_file> to pswd store
```