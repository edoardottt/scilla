#!/bin/bash

# =======================
# Scilla - Information Gathering Tool
# =======================

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see http://www.gnu.org/licenses/.

# 	@Repository:  https://github.com/edoardottt/scilla

# 	@Author:      edoardottt, https://www.edoardoottavianelli.it

# 	@License: https://github.com/edoardottt/scilla/blob/main/LICENSE



# check if HOME env var is set
if [ -z $HOME ]; then
    echo "Please set the HOME environment variable."
    exit 1
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
