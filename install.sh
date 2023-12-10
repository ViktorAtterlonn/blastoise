#!/bin/bash

download() {
  if command -v curl > /dev/null 2>&1; then
    curl -fsSL "$1"
  else
    wget -qO- "$1"
  fi
}

download_and_install() {
  if ! [ -d ~/.blastoise ]; then
    mkdir ~/.blastoise || abort "Could not create directory!"
  fi

  bin_url=`find_bin_file`

  download "$bin_url" > ~/.blastoise/blastoise  || return 1
  chmod +x ~/.blastoise/blastoise || return 1

  # add_to_path
}

find_bin_file() {
  repository="ViktorAtterlonn/blastoise"
  file_name="blastoise"

  release_info=$(curl -s "https://api.github.com/repos/$repository/releases/latest")

  # Extract the download URL
  download_url=""
  while IFS= read -r line; do
    if [[ $line =~ "browser_download_url" ]]; then
      download_url="${line#*\": \"}"
      download_url="${download_url%\"*}"
      break
    fi
  done <<< "$release_info"

  if [ -n "$download_url" ]; then
      echo $download_url
  else
      echo "File not found in the latest release."
  fi
}

add_to_path() {
  echo "export PATH=$PATH:~/.blastoise" >> ~/.zshrc
  # echo "export PATH=$PATH:$path" >> ~/.bash_profile
}

download_and_install || abort "Install Error!"