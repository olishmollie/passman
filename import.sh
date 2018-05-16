#!/bin/sh

while read -r line
do
  passman touch $line
done < "$1"
