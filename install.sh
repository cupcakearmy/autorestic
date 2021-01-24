#!/bin/bash

OUT_FILE=/usr/local/bin/autorestic

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    TYPE=linux
elif [[ "$OSTYPE" == "darwin"* ]]; then
    TYPE=macos
else
    echo "Unsupported OS"
    exit
fi

curl -s https://api.github.com/repos/cupcakearmy/autorestic/releases/latest \
| grep "browser_download_url.*_${TYPE}" \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -O ${OUT_FILE} -i -
chmod +x ${OUT_FILE}

autorestic install
autorestic --help
