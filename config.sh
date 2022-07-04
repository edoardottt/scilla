#!/bin/bash

# check if HOME env var is set
if [ -z $HOME ]; then
    echo "Please set the HOME environment variable."
fi

# check if scilla configuration folder exists 
if [[ ! -d "$HOME/.config/scilla" ]]; then
    mkdir -p $HOME/.config/scilla
fi

# check if scilla keys file exists 
if [[ ! -e "$HOME/.config/scilla/keys.yaml" ]]; then
    touch $HOME/.config/scilla/keys.yaml
    echo -n "VirusTotal: " > $HOME/.config/scilla/keys.yaml
    echo "Remember to add your keys in $HOME/.config/scilla/keys.yaml!"
fi

echo "Configuration OK."
