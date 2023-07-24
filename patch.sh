#!/usr/bin/env bash

a=$(printf "QW5vbnltZQ==" | base64 --decode)
b=$(printf "SmVzdXM=" | base64 --decode)

if [ -z $1 ]; then exit 1; fi
sed -i "s/$a/$b/gi" $1
