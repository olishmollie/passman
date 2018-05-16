#!/bin/sh
[ -f "$1" ] && rm $1
for i in $(find $HOME/.passman -type f)
do
  if [ "${i#$HOME/.passman/}" = ".fpubkey" ]
  then
    continue
  fi
  echo "${i#$HOME/.passman/} $(passman ${i#~/.passman/})" >> $1
done
