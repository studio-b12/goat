#!/bin/bash

function run_test_folder {
    for case in "$1"/*; do
        pushd "$case" &> /dev/null
        LOG=$(bash run.sh) && {
            echo -e "\e[42mSUCCESS \e[0m $case"
        } || {
            echo -e "\e[41mERROR   \e[0m $case"
            printf "%s\n" "$LOG"
            exit 1
        }
        popd &> /dev/null
    done
}

# -----------------------------------------------------------------

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

set -e

cd "$(dirname "$0")"

echo -e "\e[36mRunning cases ...\e[0m"
run_test_folder "cases"

echo -e "\n\e[36mRunning issues ...\e[0m"
run_test_folder "issues"

kill %%