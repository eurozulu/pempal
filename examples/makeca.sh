#!/bin/zsh

function makeconfig() {
  echo "Creating config file " "$1"
  cp defaultconfig "$1"
  ../bin/pp config "$1"
}

function confirm {
  read -q "response?$1? [y/N]"
  case $response in
  [yY] ) echo -e "\n"
         return;;
  [nN] ) echo -e "\naborting...\n";
         exit;;
  [*] ) echo "type y or n";;
  esac
}

CANAME="$1"
if [ -z "$CANAME" ];then
  echo "Provide an organisation name under which to create the certificate authority"
  exit
fi

CAPATH="$PWD/$CANAME"
if [ -n "$2" ];then
  CAPATH="$PWD/$2"
fi

if [ ! -d "$CAPATH" ];then
  confirm "Create directory: $CAPATH "
  mkdir -p "$CAPATH"
fi

if [ ! -f "$CAPATH/.config" ];then
  makeconfig "$CAPATH/.config"
fi

