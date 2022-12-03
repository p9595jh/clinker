#!/bin/bash

mode="$1"

if [ -z $mode ]; then

    echo -e "no parameter provided"

elif [ $mode = "start" ]; then

    . .env
    ganache \
        -p=8545 \
        --server.ws=true \
        --chain.chainId=1207 \
        -h=0.0.0.0 \
        -b=5 \
        --wallet.accounts="$PK1,90000000000000000000000" \
        --wallet.accounts="$PK2,90000000000000000000000"

elif [ $mode = "stop" ]; then

    ps=($(ps | grep 'ganache'))
    kill -9 ${ps[0]}

else

    echo -e "invalid mode: $mode (must be 'start' or 'stop')"

fi
