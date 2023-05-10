#!/bin/bash

which "echo-server" &> /dev/null || {
    go install github.com/zekroTJA/echo/cmd/echo@latest
}

which "goat" &> /dev/null || {
    task install
}

export ECHO_ADDR=localhost:8080
export ECHO_VERBOSITY=4

export GOAT_INSTANCE=http://localhost:8080

"$HOME/go/bin/echo" &> /dev/null &
sleep 1

cd $(dirname $0)

for case in cases/*; do
    pushd $case &> /dev/null
    LOG=$(bash run.sh) && {
        echo -e "\e[42mSUCCESS \e[0m $case"
    } || {
        echo -e "\e[41mERROR   \e[0m $case"
        printf "$LOG\n"
        exit 1
    }
    popd &> /dev/null
done

kill %%