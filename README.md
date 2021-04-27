<p align="center">
  <img src="https://github.com/edoardottt/images/blob/main/scilla/logo.png"><br>
  <b>🏴‍☠️ Information Gathering tool 🏴‍☠️ - DNS / Subdomains / Ports / Directories enumeration</b><br>
  <br>
  <!-- go-report-card -->
  <a href="https://goreportcard.com/report/github.com/edoardottt/scilla">
    <img src="https://goreportcard.com/badge/github.com/edoardottt/scilla" alt="go-report-card" />
  </a>
  <!-- workflows -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/scilla/workflows/Go/badge.svg?branch=master" alt="workflows" />
  </a>
  <!-- ubuntu-build -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/images/blob/main/scilla/ubuntu-build.svg" alt="ubuntu-build" />
  </a>
  <!-- win10-build -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/images/blob/main/scilla/win10.svg" alt="win10-build" />
  </a>
  <!-- pr-welcome -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/images/blob/main/scilla/pr-welcome.svg" alt="pr-welcome" />
  </a>

  <br>
  
  <!-- mainteinance -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/images/blob/main/scilla/maintained-yes.svg" alt="Mainteinance yes" />
  </a>
  <!-- ask-me-anything -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/images/blob/main/scilla/ask-me-anything.svg" alt="ask me anything" />
  </a>
  <!-- gobadge -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/images/blob/main/scilla/gobadge" alt="gobadge" />
  </a>
  <!-- license GPLv3.0 -->
  <a href="https://github.com/edoardottt/scilla/blob/master/LICENSE">
    <img src="https://github.com/edoardottt/images/blob/main/scilla/license-GPL3.svg" alt="license-GPL3" />
  </a>
  <br>
  <sub>
    Coded with 💙 by edoardottt.
  </sub>
  <br>
  <!--Tweet button-->
  <a href="https://twitter.com/intent/tweet?url=https%3A%2F%2Fgithub.com%2Fedoardottt%2Fscilla%20&text=Information%20Gathering%20tool%21&hashtags=pentesting%2Clinux%2Cgolang%2Cnetwork" target="_blank">Share on Twitter!
  </a>
</p>
<p align="center">
  <a href="#preview-bar_chart">Preview</a> •
  <a href="#installation-">Install</a> •
  <a href="#get-started-">Get Started</a> •
  <a href="#examples-">Examples</a> •
  <a href="#contributing-">Contributing</a>
</p>

Preview :bar_chart:
----------

[![asciicast](https://asciinema.org/a/406707.svg)](https://asciinema.org/a/406707)

Installation 📡
----------

- **Linux**

  - `git clone https://github.com/edoardottt/scilla.git`
  - `cd scilla`
  - `go get`
  - `make linux` (to install)
  - `make unlinux` (to uninstall)

- **Windows** (executable works only in scilla folder. [Alias?](https://github.com/edoardottt/scilla/issues/10))

  - `git clone https://github.com/edoardottt/scilla.git`
  - `cd scilla`
  - `go get`
  - `.\make.bat windows` (to install)    
  - `.\make.bat unwindows` (to uninstall)

Get Started 🎉
----------

`scilla help` prints the help in the command line.

```
usage: scilla subcommand { options }

   Available subcommands:
       - dns [-o output-format]
             [-plain Print only results]
             -target <target (URL/IP)> REQUIRED
       - port [-p <start-end> or ports divided by comma]
              [-o output-format]
              [-common scan common ports]
              [-plain Print only results]
              -target <target (URL/IP)> REQUIRED
       - subdomain [-w wordlist]
                   [-o output-format]
                   [-i ignore status codes]
                   [-c use also a web crawler]
                   [-db use also a public database]
                   [-plain Print only results]
                   -target <target (URL)> REQUIRED
       - dir [-w wordlist]
             [-o output-format]
             [-i ignore status codes]
             [-c use also a web crawler]
             [-plain Print only results]
             -target <target (URL)> REQUIRED
       - report [-p <start-end> or ports divided by comma]
                [-ws subdomains wordlist]
                [-wd directories wordlist]
                [-o output-format]
                [-id ignore status codes in directories scanning]
                [-is ignore status codes in subdomains scanning]
                [-cd use also a web crawler for directories scanning]
                [-cs use also a web crawler for subdomains scanning]
                [-db use also a public database for subdomains scanning]
                [-common scan common ports]
                -target <target (URL/IP)> REQUIRED
       - help
       - examples
```


Examples 💡
----------

- DNS enumeration:
    
    - `scilla dns -target target.domain`
    - `scilla dns -o txt -target target.domain`
    - `scilla dns -o html -target target.domain`
    - `scilla dns -plain -target target.domain`

- Subdomains enumeration:

    - `scilla subdomain -target target.domain`
    - `scilla subdomain -w wordlist.txt -target target.domain`
    - `scilla subdomain -o txt -target target.domain`
    - `scilla subdomain -o html -target target.domain`
    - `scilla subdomain -i 400 -target target.domain`
    - `scilla subdomain -i 4** -target target.domain`
    - `scilla subdomain -c -target target.domain`
    - `scilla subdomain -db -target target.domain`
    - `scilla subdomain -plain -target target.domain`

- Directories enumeration:

    - `scilla dir -target target.domain`
    - `scilla dir -w wordlist.txt -target target.domain`
    - `scilla dir -o txt -target target.domain`
    - `scilla dir -o html -target target.domain`
    - `scilla dir -i 500,401 -target target.domain`
    - `scilla dir -i 5**,401 -target target.domain`
    - `scilla dir -c -target target.domain`
    - `scilla dir -plain -target target.domain`

- Ports enumeration:
      
    - Default (all ports, so 1-65635) `scilla port -target target.domain`
    - Specifying ports range `scilla port -p 20-90 -target target.domain`
    - Specifying starting port (until the last one) `scilla port -p 20- -target target.domain`
    - Specifying ending port (from the first one) `scilla port -p -90 -target target.domain`
    - Specifying single port `scilla port -p 80 -target target.domain`
    - Specifying output format (txt)`scilla port -o txt -target target.domain`
    - Specifying output format (html)`scilla port -o html -target target.domain`
    - Specifying multiple ports `scilla port -p 21,25,80 -target target.domain`
    - Specifying common ports `scilla port -common -target target.domain`
    - Print only results `scilla port -plain -target target.domain`

- Full report:
      
    - Default (all ports, so 1-65635) `scilla report -target target.domain`
    - Specifying ports range `scilla report -p 20-90 -target target.domain`
    - Specifying starting port (until the last one) `scilla report -p 20- -target target.domain`
    - Specifying ending port (from the first one) `scilla report -p -90 -target target.domain`
    - Specifying single port `scilla report -p 80 -target target.domain`
    - Specifying output format (txt)`scilla report -o txt -target target.domain`
    - Specifying output format (html)`scilla report -o html -target target.domain`
    - Specifying directories wordlist `scilla report -wd dirs.txt -target target.domain`
    - Specifying subdomains wordlist `scilla report -ws subdomains.txt -target target.domain`
    - Specifying status codes to be ignored in directories scanning `scilla report -id 500,501,502 -target target.domain`
    - Specifying status codes to be ignored in subdomains scanning `scilla report -is 500,501,502 -target target.domain`
    - Specifying status codes classes to be ignored in directories scanning `scilla report -id 5**,4** -target target.domain`
    - Specifying status codes classes to be ignored in subdomains scanning `scilla report -is 5**,4** -target target.domain`
    - Use also a web crawler for directories enumeration `scilla report -cd -target target.domain`
    - Use also a web crawler for subdomains enumeration `scilla report -cs -target target.domain`
    - Use also a public database for subdomains enumeration `scilla report -db -target target.domain`
    - Specifying multiple ports `scilla report -p 21,25,80 -target target.domain`
    - Specifying common ports `scilla report -common -target target.domain`

Contributing 🛠
-------
<!--
[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/0)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/0)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/1)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/1)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/2)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/2)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/3)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/3)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/4)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/4)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/5)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/5)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/6)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/6)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/7)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/7)
-->
Just open an issue/pull request. See also [CONTRIBUTING.md](https://github.com/edoardottt/scilla/blob/master/CONTRIBUTING.md) and [CODE OF CONDUCT.md](https://github.com/edoardottt/scilla/blob/master/CODE_OF_CONDUCT.md)

**Help me building this!**

A special thanks to [danielmiessler](https://github.com/danielmiessler), using those lists; and to [sonarSearch](https://github.com/cgboal/sonarsearch).

**To do:**

  - [ ] Tests (😂)
  
  - [ ] Tor support
  
  - [ ] Proxy support
    
  - [ ] JSON output
  
  - [ ] XML output

  - [x] Plain output (print only results)

  - [x] Scan only common ports

  - [x] Add option to use a public database of known subdomains

  - [x] Recursive Web crawling for subdomains and directories

  - [x] Check input and if it's an IP try to change to hostname when dns or subdomain is active
  
  - [x] Ignore responses by status codes (partially done, to do with `*`, e.g. `-i 4**`)
  
  - [x] HTML output
  
  - [x] Build an Input Struct and use it as parameter

  - [x] Output color
  
  - [x] Subdomains enumeration
  
  - [x] DNS enumeration
 
  - [x] Subdomains enumeration

  - [x] Port enumeration

  - [x] Directories enumeration
  
  - [x] TXT output
  
License 📝
-------

This repository is under [GNU General Public License v3.0](https://github.com/edoardottt/scilla/blob/main/LICENSE).  
[edoardoottavianelli.it](https://www.edoardoottavianelli.it) to contact me.
