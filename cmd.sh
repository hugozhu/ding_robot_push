#!/bin/bash
line=$@
line=$(echo -e "${line}" | sed -e 's/^[[:space:]]*//')

if [[ "$line" == /* ]]; then
    cmd1=`echo $line | cut -d ' ' -f1`
    cmd=${cmd1:1:10000}
    if [[ "ip" == $cmd ]]; then 
        ip=$(printf "\n" | /bin/nc ns1.dnspod.net 6666)
        push3.sh "$ip"
    fi

    if [[ "echo" == $cmd ]]; then
        t=`echo $line | cut -d ' ' -f2`
        push3.sh "$t"
    fi
fi