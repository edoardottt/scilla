# Scilla
<p align="center">
  <!-- logo -->
  <b>🏴‍☠️ Information Gathering tool 🏴‍☠️ - DNS / Subdomain / Ports / Directories enumeration</b><br>
    <sub>
    Coded with 💙 by edoardottt.
  </sub>
  <br>
  <!--Tweet button-->
  <a href="https://twitter.com/intent/tweet?url=https%3A%2F%2Fgithub.com%2Fedoardottt%2Fscilla%20&text=Information%20Gathering%20tool%21&hashtags=pentesting%2Clinux%2Cgolang%2Cnetwork" target="_blank">Share on Twitter!
  </a>
  <br>
  <!-- mainteinance -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/scilla/blob/master/images/maintained-yes.svg" alt="Mainteinance yes" />
  </a>
  <!-- pr-welcome -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/scilla/blob/master/images/pr-welcome.svg" alt="pr-welcome" />
  </a>
  <!-- ask-me-anything -->
  <a href="https://edoardoottavianelli.it">
    <img src="https://github.com/edoardottt/scilla/blob/master/images/ask-me-anything.svg" alt="ask me anything" />
  </a>
    <!-- workflows -->
      <a href="https://edoardoottavianelli.it">
        <img src="https://github.com/edoardottt/scilla/workflows/Go/badge.svg?branch=master" alt="workflows" />
      </a>
    <br>
    <!-- ubuntu-build -->
      <a href="https://edoardoottavianelli.it">
        <img src="https://github.com/edoardottt/scilla/blob/master/images/ubuntu-build.svg" alt="ubuntu-build" />
      </a>
    <!-- go-report-card -->
      <a href="https://goreportcard.com/report/github.com/edoardottt/scilla">
        <img src="https://goreportcard.com/badge/github.com/edoardottt/scilla" alt="go-report-card" />
      </a>
    <!-- gobadge -->
      <a href="https://edoardoottavianelli.it">
        <img src="https://github.com/edoardottt/scilla/blob/master/images/gobadge" alt="gobadge" />
      </a>
    <!-- license GPLv3.0 -->
      <a href="https://github.com/edoardottt/scilla/blob/master/LICENSE">
        <img src="https://github.com/edoardottt/scilla/blob/master/images/license-GPL3.svg" alt="license-GPL3" />
      </a>
</p>

Example :bar_chart:
----------

![Example](https://github.com/edoardottt/scilla/blob/master/images/scilla.gif)

Installation 📡
----------

First of all, clone the repo locally

`git clone https://github.com/edoardottt/scilla.git`

Scilla has external dependencies, so they need to be pulled in:

`go get`

**Working on installation...** [See the open issue](https://github.com/edoardottt/scilla/issues/4).

For now you can run it inside the scilla folder with `go run scilla.go ...`  

Too late.. : see [this](https://www.maketecheasier.com/make-scripts-executable-everywhere-linux/)

Then use the build scripts:

- `make windows` builds 32 and 64 bit binaries for Windows, and writes them to the build subfolder.

- `make linux` builds 32 and 64 bit binaries for Linux, and writes them to the build subfolder.

- `make unlinux` Removes binaries.

- `make fmt` run the golang formatter.

- `make update` Update.

- `make remod` Remod.

- `make test` runs the tests.

- `make clean` clears out the build subfolder.


Get Started 🎉
----------

`scilla help` prints the help in the command line.

    usage: scilla [subcommand] { options }
    
        Available subcommands:
            - dns { -target <target (URL)> REQUIRED}
            - subdomain { [-w wordlist] -target <target (URL)> REQUIRED}
            - port { [-p <start-end>] -target <target (URL/IP)> REQUIRED}
            - dir { [-w wordlist] -target <target (URL/IP)> REQUIRED}
            - report { [-p <start-end>] -target <target (URL/IP)> REQUIRED}
            - help


Examples 💡
----------

- DNS enumeration:
    
    - `scilla dns -target target.domain`

- Subdomains enumeration:

    - `scilla subdomain -target target.domain`

    - `scilla subdomain -w wordlist.txt -target target.domain`

- Directories enumeration:

    - `scilla dir -target target.domain`

    - `scilla dir -w wordlist.txt -target target.domain`

- Ports enumeration:
      
    - Default (all ports, so 1-65635) `scilla port -target target.domain`

    - Specifying ports range `scilla port -p 20-90 -target target.domain`

    - Specifying starting port (until the last one) `scilla port -p 20- -target target.domain`

    - Specifying ending port (from the first one) `scilla port -p -90 -target target.domain`

    - Specifying single port `scilla port -p 80 -target target.domain`

- Full report:
      
    - Default (all ports, so 1-65635) `scilla report -target target.domain`

    - Specifying ports range `scilla report -p 20-90 -target target.domain`

    - Specifying starting port (until the last one) `scilla report -p 20- -target target.domain`

    - Specifying ending port (from the first one) `scilla report -p -90 -target target.domain`

    - Specifying single port `scilla report -p 80 -target target.domain`

Contributing 🛠
-------

[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/0)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/0)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/1)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/1)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/2)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/2)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/3)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/3)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/4)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/4)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/5)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/5)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/6)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/6)[![](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/images/7)](https://sourcerer.io/fame/edoardottt/edoardottt/scilla/links/7)

Just open an issue/pull request. See also [CONTRIBUTING.md](https://github.com/edoardottt/scilla/blob/master/CONTRIBUTING.md) and [CODE OF CONDUCT.md](https://github.com/edoardottt/scilla/blob/master/CODE_OF_CONDUCT.md)

**Help me building this!**

A special thanks to [danielmiessler](https://github.com/danielmiessler). Now using one of those lists.

**To do:**

  - [ ] Test the functions built
  
  - [x] Subdomains enumeration
  
  - [x] DNS enumeration
 
  - [x] Subdomains enumeration

  - [x] Port enumeration

  - [x] Directories enumeration
  
  - [ ] Print the progress percentage value when CR is pressed (not in output doc)
  
  - [x] Build an Input Struct and use it as parameter

  - [x] Output color
  
  - [ ] Check input and if it's an IP try to change to hostname when dns or subdomain is active
  
  - [ ] JSON report output
  
  - [ ] PDF report output
  
  - [ ] XML report output
  
  - [ ] (report mode) In all the subdomains found enumerates ports???
  
  - [ ] Tor support
  
  - [ ] Proxy support


If you liked it drop a :star:
-------

https://www.edoardoottavianelli.it for contact me.


  
                                                                    Edoardo Ottavianelli
