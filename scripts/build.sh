#!/bin/bash

set -e

version=$1

GROUP=kingsoftcloud
SHORT_NAME=ksyun
PLUGIN_NAME=packer-plugin-${SHORT_NAME}

# Detech current os category
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     OS_TYPE=linux;;
    Darwin*)    OS_TYPE=darwin;;
    CYGWIN*)    OS_TYPE=windows;;
    MINGW*)     OS_TYPE=windows;;
    *)          OS_TYPE="UNKNOWN:${unameOut}"
esac

unameOut="$(uname -m)"
case "${unameOut}" in
    x86_64*)     OS_ARCH=amd64;;
    arm64*)      OS_ARCH=arm64;;
    *)          OS_ARCH="UNKNOWN:${unameOut}"
esac

FULL_PLUGIN_NAME=${PLUGIN_NAME}_${version}_x5.0_${OS_TYPE}_${OS_ARCH}
echo "${FULL_PLUGIN_NAME} is deteched."
echo "Compiling ..."

plugin_path=$HOME/.packer.d/plugins/github.com/${GROUP}/${SHORT_NAME}

echo $plugin_path/$FULL_PLUGIN_NAME

if [ $OS_TYPE == "linux" -o $OS_TYPE == "darwin" ]; then
	GOOS=$OS_TYPE GOARCH=$OS_ARCH go build -ldflags "-X main.version=$version" -o bin/${PLUGIN_NAME}
	chmod +x bin/${PLUGIN_NAME}
    mkdir -p $plugin_path
    mv bin/${PLUGIN_NAME} $plugin_path/${FULL_PLUGIN_NAME}
    shasum -a 256 $plugin_path/${FULL_PLUGIN_NAME} > $plugin_path/${FULL_PLUGIN_NAME}_SHA256SUM
elif [ $OS_TYPE == "Windows" ]; then
	GOOS=$OS_TYPE GOARCH=$OS_ARCH  go build -ldflags "-X main.version=$version" -o bin/${PLUGIN_NAME}
	chmod +x bin/${PLUGIN_NAME}.exe
    mkdir -p $plugin_path_win
    mv bin/${PLUGIN_NAME}.exe $plugin_path_win/${PLUGIN_NAME}.exe
else
    echo "Invalid OS"
    exit 1
fi