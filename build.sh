#!/usr/bin/bash

o="../bin/scoreboard"
if [ $# -ge 1 ]; then
    o="$1"
fi

which packr2 &> /dev/null
if [ $? -ne 0 ]; then
    printf "Installing packr..\n"
    bash -c "cd scoreboard; go get -u github.com/gobuffalo/packr/v2/packr2"
fi

printf "Building...\n"
bash -c "cd scoreboard; packr2; go build -o \"$o\" cmd/scoreboard/main.go; packr2 clean"
printf "Done!\n"
