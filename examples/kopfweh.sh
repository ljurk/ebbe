#!/bin/zsh
ebbe merge --commands -i \
    "ebbe color -c 000000 -c ffffff  -w 1024 -h 256" \
    "ebbe color -v -y 256 -c 000000 -c ffffff -w 1024 -h 256" \
    -i "ebbe text --text acab --fontsize 30" \
    | ebbe send --host 10.13.12.196:1337 --input - -c 64 -p 2048
