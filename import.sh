#!/bin/sh

# Alternate implementation of `passman import`

while read -r line
do
  passman touch $line
done < "$1"
