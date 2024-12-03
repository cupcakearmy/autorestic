#!/bin/bash
# Check if the script is running as root
# Initialize the elevated variable to ensure we can run commands as root
if [ "$EUID" -eq 0 ]; then
    elevated="bash -c"  # No need for sudo or su if running as root
elif command -v sudo &> /dev/null; then
    elevated="sudo bash -c"  # Use sudo if available
elif command -v su &> /dev/null; then
    elevated="su -c"  # Use su if sudo isn't available
else
    echo "Error: Neither sudo nor su is available. Unable to run command as root."
    exit 1
fi
set -e -o pipefail
shopt -s nocaseglob

OUT_FILE=/usr/local/bin/autorestic
TMP_FILE=/tmp/autorestic

# Type
NATIVE_OS=$(uname | tr '[:upper:]' '[:lower:]')
case $NATIVE_OS in
  *linux*)
    OS=linux;;
  *darwin*)
    OS=darwin;;
  *freebsd*)
    OS=freebsd;;
  *)
    echo "Could not determine OS automatically, please check the release page"\
    "manually: https://github.com/cupcakearmy/autorestic/releases"
    exit 1
    ;;
esac
echo $OS

NATIVE_ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')
case $NATIVE_ARCH in
  *x86_64*|*amd64*) ARCH=amd64 ;;
  *arm64*|*aarch64*) ARCH=arm64 ;;
  *x86*) ARCH=386 ;;
  *armv7*) ARCH=arm ;;
        *)
          echo "Could not determine Architecure automatically, please check the"\
          "release page manually: https://github.com/cupcakearmy/autorestic/releases"
          exit 1
          ;;
esac
echo $ARCH

if ! command -v bzip2 &>/dev/null; then
    echo "Missing bzip2 command. Please install the bzip2 package for your system."
    exit 1
fi

wget -qO - https://api.github.com/repos/cupcakearmy/autorestic/releases/latest \
| grep "browser_download_url.*_${OS}_${ARCH}" \
| xargs | cut -d ' ' -f 2 \
| wget -O "${TMP_FILE}.bz2" -i -
$elevated "bzip2 -cd ${TMP_FILE}.bz2 > ${OUT_FILE}"
$elevated "chmod +x ${OUT_FILE}"
$elevated "rm ${TMP_FILE}.bz2"

autorestic install
echo "Successfully installed autorestic under ${OUT_FILE}"
