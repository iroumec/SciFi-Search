#!/bin/bash
exec 3<>/dev/tcp/127.0.0.1/$1
echo -e "GET /hello HTTP/1.1\r\nhost: 127.0.0.1:$1\r\nConnection: close\r\n\r\n" >&3
cat <&3 | grep "Hello"