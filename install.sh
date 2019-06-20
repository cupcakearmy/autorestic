#!/bin/sh
OUT_FILE=/usr/local/bin/autorestic

curl -s https://api.github.com/repos/cupcakearmy/autorestic/releases/latest \
| grep "browser_download_url.*-macos" \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -q -O $OUT_FILE -i -
chmod +x $OUT_FILE