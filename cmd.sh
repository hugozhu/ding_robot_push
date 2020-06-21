#!/bin/bash
line=$@
line=$(echo "${line}" | sed -e 's/^[[:space:]]*//')
if [[ "$line" == /* ]]; then
    cmd1=`echo $line | cut -d ' ' -f1`
    cmd=${cmd1:1:10000}

    if [[ "help" == $cmd ]]; then
        help=$(cat << END
        USAGE:
        /ip:    显示当前公网IP 
        /echo:  回显第一个单词 
        /temp:  显示树莓派CPU温度 
        /sh:    显示市场股价 
        /cny:   显示市场汇率
        /cam:   获取视像头拍照
END
)
       push3.sh "$help"
    fi
    
    if [[ "cam" == $cmd ]]; then
        echo "screenshot capturing..."
        /home/hugo/Projects/ding_robot_push/cam.sh
    fi

    if [[ "ip" == $cmd ]]; then 
        ip=$(printf "\n" | /bin/nc ns1.dnspod.net 6666)
        push3.sh "$ip"
    fi

    if [[ "echo" == $cmd ]]; then
        t=`echo $line | cut -d ' ' -f2`
        push3.sh "$t"
    fi

    if [[ "temp" == $cmd ]]; then
        t=`/usr/bin/vcgencmd measure_temp`
        push3.sh "$t"
    fi

    if [[ "sh" == $cmd ]]; then
        t=`/home/hugo/bin/phantomjs2 /home/hugo/Projects/phantomjs/sh00001.js > /dev/null`
    fi

    if [[ "cny" == $cmd ]]; then
        /home/hugo/Projects/phantomjs/safe.gov.sh
    fi
fi
