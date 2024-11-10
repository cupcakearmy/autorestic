#!/bin/bash
set -e -o pipefail
shopt -s nocaseglob

OUT_FILE=/usr/local/bin/autorestic

# Type
NATIVE_OS=$(uname | tr '[:upper:]' '[:lower:]')
if [[ $NATIVE_OS == *"linux"* ]]; then
    OS=linux
elif [[ $NATIVE_OS == *"darwin"* ]]; then
    OS=darwin
elif [[ $NATIVE_OS == *"freebsd"* ]]; then
    OS=freebsd
else
    echo "Could not determine OS automatically, please check the release page manually: https://github.com/cupcakearmy/autorestic/releases"
    exit 1
fi
echo $OS

NATIVE_ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')
if [[ $NATIVE_ARCH == *"x86_64"* || $NATIVE_ARCH == *"amd64"* ]]; then
    ARCH=amd64
elif [[ $NATIVE_ARCH == *"arm64"* || $NATIVE_ARCH == *"aarch64"* ]]; then
    ARCH=arm64
elif [[ $NATIVE_ARCH == *"x86"* ]]; then
    ARCH=386
elif [[ $NATIVE_ARCH == *"armv7"* ]]; then
    ARCH=arm
else
    echo "Could not determine Architecure automatically, please check the release page manually: https://github.com/cupcakearmy/autorestic/releases"
    exit 1
fi
echo $ARCH

if ! command -v bzip2 &>/dev/null; then
    echo "Missing bzip2 command. Please install the bzip2 package for your system."
    exit 1
fi

wget -qO - https://api.github.com/repos/cupcakearmy/autorestic/releases/latest \
| grep "browser_download_url.*_${OS}_${ARCH}" \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -O "${OUT_FILE}.bz2" -i -
bzip2 -fd "${OUT_FILE}.bz2"
chmod +x ${OUT_FILE}

autorestic install
echo "Successfully installed autorestic"
