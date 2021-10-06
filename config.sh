#!/bin/bash

if [[ ! -d "~/.config/scilla" ]]; then
    mkdir -p ~/.config/scilla
fi
if [[ ! -f "~/.config/scilla/keys.yaml" ]]; then
    touch ~/.config/scilla/keys.yaml
    echo -n "Spyse: 
VirusTotal: " > ~/.config/scilla/keys.yaml
fi
echo "Configuration OK."