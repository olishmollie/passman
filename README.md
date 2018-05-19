# passman

#### Installation
```
$ go get github.com/olishmollie/passman
$ go install github.com/olishmollie/passman
```

#### Getting started
Passman stores encrypted passwords in a directory in your home folder. To get started, run `passman init`, which will create the password store at `~/.passman` and generate an encryption key. This key is stored in `~/.passman/.key`.

#### Tab Completion
Passman comes bundled with a bash completion script. To use it, source passman-completion.bash from your bashrc or bash_profile.

#### Usage
```
usage: passman [command] [args...] [opts...]

commands: dump edit generate import init rm touch

passman - prints a tree of pswds in store

passman [opts...] <pswd_file> - prints unencrypted pswd
    -c, --copy
        copy password to clipboard

dump - prints unencrypted pswds to stdout

edit <pswd_file> - edit pswd in editor set to $VISUAL

generate [opts...] - generates a random pswd
    -c, --copy
        copies unencrypted pswd to clipboard
    -l, --len int 
        specifies length of generated pswd
    -n, --nosym 
        generate a password with no symbols

import <infile> - imports passwords from infile.
    NOTE: infile must be in the following format:
        website/username secret_password
        Category/anothersite/username another_password
        etc.

init - create pswd store if it doesn't exist, generate encryption key

lock - encrypts and dumps all passwords into one file

rm <pswd_file> - remove <pswd_file> from pswd store

touch <pswd_file> - add <pswd_file> to pswd store

unlock - undoes lock operation
```