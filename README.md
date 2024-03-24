<p align="center">
  <img src="https://github.com/edoardottt/images/blob/main/scilla/logo.png"><br>
  <b>🏴‍☠️ Information Gathering tool 🏴‍☠️ - DNS / Subdomains / Ports / Directories enumeration</b><br>
  <br>
  <!-- go-report-card -->
  <a href="https://goreportcard.com/report/github.com/edoardottt/scilla">
    <img src="https://goreportcard.com/badge/github.com/edoardottt/scilla" alt="go-report-card" />
  </a>
  <!-- workflows -->
  <a href="https://github.com/edoardottt/scilla/actions">
    <img src="https://github.com/edoardottt/scilla/workflows/Go/badge.svg?branch=main" alt="workflows" />
  </a>
  <br>
  <sub>
    Coded with 💙 by edoardottt
  </sub>
  <br>
  <!--Tweet button-->
  <a href="https://twitter.com/intent/tweet?url=https%3A%2F%2Fgithub.com%2Fedoardottt%2Fscilla%20&text=Information%20Gathering%20tool%21&hashtags=pentesting%2Clinux%2Cgolang%2Cnetwork" target="_blank">Share on Twitter!
  </a>
</p>
<p align="center">
  <a href="#installation-">Install</a> •
  <a href="#get-started-">Get Started</a> •
  <a href="#examples-">Examples</a> •
  <a href="#changelog-">Changelog</a> •
  <a href="#contributing-">Contributing</a> •
  <a href="#license-">License</a>
</p>

<!--[![asciicast](https://asciinema.org/a/406707.svg)](https://asciinema.org/a/406707)-->

<p align="center">
  <img src="https://github.com/edoardottt/images/blob/main/scilla/scilla.gif">
</p>

Installation 📡
----------

### Building from source

You need [Go](https://golang.org/).

- **Linux**

  - `git clone https://github.com/edoardottt/scilla.git`
  - `cd scilla`
  - `make linux` (to install)
  - Edit the `~/.config/scilla/keys.yaml` file if you want to use API keys
  - `make unlinux` (to uninstall)

- **Windows** (executable works only in scilla folder. [Alias?](https://github.com/edoardottt/scilla/issues/10))

  - `git clone https://github.com/edoardottt/scilla.git`
  - `cd scilla`
  - `.\make.bat windows` (to install)
  - Create a `keys.yaml` file  if you want to use api keys
  - `.\make.bat unwindows` (to uninstall)

### Using Docker

```shell
docker build -t scilla .
docker run scilla help
```

Get Started 🎉
----------

`scilla help` prints the help in the command line.

```
usage: scilla subcommand { options }

   Available subcommands:
       - dns [-oj JSON output file]
             [-oh HTML output file]
             [-ot TXT output file]
             [-plain Print only results]
             -target <target (URL/IP)> REQUIRED
       - port [-p <start-end> or ports divided by comma]
              [-oj JSON output file]
              [-oh HTML output file]
              [-ot TXT output file]
              [-common scan common ports]
              [-plain Print only results]
              -target <target (URL/IP)> REQUIRED
       - subdomain [-w wordlist]
                   [-oj JSON output file]
                   [-oh HTML output file]
                   [-ot TXT output file]
                   [-i ignore status codes]
                   [-c use also a web crawler]
                   [-db use also a public database]
                   [-plain Print only results]
                   [-db -no-check Don't check status codes for subdomains]
                   [-db -vt Use VirusTotal as subdomains source]
                   [-db -bw Use BuiltWith as subdomains source]
                   [-ua Set the User Agent]
                   [-rua Generate a random user agent for each request]
                   [-dns Set DNS IP to resolve the subdomains]
                   [-alive Check also if the subdomains are alive]
                   -target <target (URL)> REQUIRED
       - dir [-w wordlist]
             [-oj JSON output file]
             [-oh HTML output file]
             [-ot TXT output file]
             [-i ignore status codes]
             [-c use also a web crawler]
             [-plain Print only results]
             [-nr No follow redirects]
             [-ua Set the User Agent]
             [-rua Generate a random user agent for each request]
             -target <target (URL/IP)> REQUIRED
       - report [-p <start-end> or ports divided by comma]
                [-ws subdomains wordlist]
                [-wd directories wordlist]
                [-oj JSON output file]
                [-oh HTML output file]
                [-ot TXT output file]
                [-id ignore status codes in directories scanning]
                [-is ignore status codes in subdomains scanning]
                [-cd use also a web crawler for directories scanning]
                [-cs use also a web crawler for subdomains scanning]
                [-db use also a public database for subdomains scanning]
                [-common scan common ports]
                [-nr No follow redirects]
                [-db -vt Use VirusTotal as subdomains source]
                [-ua Set the User Agent]
                [-rua Generate a random user agent for each request]
                [-dns Set DNS IP to resolve the subdomains]
                [-alive Check also if the subdomains are alive]
                -target <target (URL)> REQUIRED
       - help
       - examples
```

Examples 💡
----------

- DNS enumeration:

  - `scilla dns -target target.domain`
  - `scilla dns -oj output -target target.domain`
  - `scilla dns -oh output -target target.domain`
  - `scilla dns -ot output -target target.domain`
  - `scilla dns -plain -target target.domain`

- Subdomains enumeration:

  - `scilla subdomain -target target.domain`
  - `scilla subdomain -w wordlist.txt -target target.domain`
  - `scilla subdomain -oj output -target target.domain`
  - `scilla subdomain -oh output -target target.domain`
  - `scilla subdomain -ot output -target target.domain`
  - `scilla subdomain -i 400 -target target.domain`
  - `scilla subdomain -i 4** -target target.domain`
  - `scilla subdomain -c -target target.domain`
  - `scilla subdomain -db -target target.domain`
  - `scilla subdomain -plain -target target.domain`
  - `scilla subdomain -db -no-check -target target.domain`
  - `scilla subdomain -db -vt -target target.domain`
  - `scilla subdomain -db -bw -target target.domain`
  - `scilla subdomain -ua "CustomUA" -target target.domain`
  - `scilla subdomain -rua -target target.domain`
  - `scilla subdomain -dns 8.8.8.8 -target target.domain`
  - `scilla subdomain -alive -target target.domain`

- Directories enumeration:

  - `scilla dir -target target.domain`
  - `scilla dir -w wordlist.txt -target target.domain`
  - `scilla dir -oj output -target target.domain`
  - `scilla dir -oh output -target target.domain`
  - `scilla dir -ot output -target target.domain`
  - `scilla dir -i 500,401 -target target.domain`
  - `scilla dir -i 5**,401 -target target.domain`
  - `scilla dir -c -target target.domain`
  - `scilla dir -plain -target target.domain`
  - `scilla dir -nr -target target.domain`
  - `scilla dir -ua "CustomUA" -target target.domain`
  - `scilla dir -rua -target target.domain`

- Ports enumeration:

  - Default (all ports, so 1-65635) `scilla port -target target.domain`
  - Specifying ports range `scilla port -p 20-90 -target target.domain`
  - Specifying starting port (until the last one) `scilla port -p 20- -target target.domain`
  - Specifying ending port (from the first one) `scilla port -p -90 -target target.domain`
  - Specifying single port `scilla port -p 80 -target target.domain`
  - Specifying output format (json)`scilla port -oj output -target target.domain`
  - Specifying output format (html)`scilla port -oh output -target target.domain`
  - Specifying output format (txt)`scilla port -ot output -target target.domain`
  - Specifying multiple ports `scilla port -p 21,25,80 -target target.domain`
  - Specifying common ports `scilla port -common -target target.domain`
  - Print only results `scilla port -plain -target target.domain`

- Full report:

  - Default (all ports, so 1-65635) `scilla report -target target.domain`
  - Specifying ports range `scilla report -p 20-90 -target target.domain`
  - Specifying starting port (until the last one) `scilla report -p 20- -target target.domain`
  - Specifying ending port (from the first one) `scilla report -p -90 -target target.domain`
  - Specifying single port `scilla report -p 80 -target target.domain`
  - Specifying output format (json)`scilla report -oj output -target target.domain`
  - Specifying output format (html)`scilla report -oh output -target target.domain`
  - Specifying output format (txt)`scilla report -ot output -target target.domain`
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
  - No follow redirects `scilla report -nr -target target.domain`
  - Use VirusTotal as subdomains source `scilla report -db -vt -target target.domain`
  - Set the User Agent `scilla report -ua "CustomUA" -target target.domain`
  - Generate a random user agent for each request `scilla report -rua -target target.domain`
  - Set DNS IP to resolve the subdomains `scilla report -dns 8.8.8.8 -target target.domain`
  - Check also if the subdomains are alive `scilla report -alive -target target.domain`

Changelog 📌
-------

Detailed changes for each release are documented in the [release notes](https://github.com/edoardottt/scilla/releases).

Contributing 🛠
-------

Just open an [issue](https://github.com/edoardottt/scilla/issues) / [pull request](https://github.com/edoardottt/scilla/pulls).

Before opening a pull request, download [golangci-lint](https://golangci-lint.run/usage/install/) and run

```bash
golangci-lint run
```

If there aren't errors, go ahead :)

**To do:**

- [ ] Add more tests
  
- [ ] Tor support
  
- [ ] Proxy support

In the news 📰
-------

- [Kali Linux Tutorials](https://kalilinuxtutorials.com/scilla/)
- [GeeksForGeeks.org](https://www.geeksforgeeks.org/scilla-information-gathering-dns-subdomain-port-enumeration/)
- [Brisk Infosec](https://www.briskinfosec.com/tooloftheday/toolofthedaydetail/Scilla)
- [Kalitut](https://kalitut.com/scilla-nformation-gathering-tool/)
  
License 📝
-------

This repository is under [GNU General Public License v3.0](https://github.com/edoardottt/scilla/blob/main/LICENSE).  
[edoardoottavianelli.it](https://www.edoardoottavianelli.it) to contact me.
