# Chhoto URL CLI
This is a companion CLI tool for [Chhoto URL](https://github.com/SinTan1729/chhoto-url). It provides simple commands for interacting with
a Chhoto URL server.

All the functions of this tool can be replicated using `curl`, but this is supposed to be easier to use. The server needs to have an API key
in order to be accessible using this CLI tool. Password based usage is supported, but discouraged.

# Installation
## Installation from source
1. Clone the repo.
```
git clone https://github.com/SinTan1729/chhoto-url-cli
```
2. Install.
```
cd chhoto-url-cli
make install
```
3. You can uninstall by running `make uninstall`.

## Installation from AUR
Use the AUR package [`chhoto-url-cli-bin`](https://aur.archlinux.org/packages/chhoto-url-cli-bin).

## Installation from LURE
This should (at least in theory) work for every distro, and should be similar to AUR in terms of experience.
1. Install `LURE` from [lure.sh](https://lure.sh).
2. Add my personal repo to it.
```
lure addrepo -n SinTan1729 -u https://github.com/SinTan1729/lure-repo
```
3. Install `chhoto-url-cli`
```
lure in chhoto-url-cli
```

# Usage
By default, config will be loaded from $XDG_CONFIG_HOME/chhoto/config.json
But these can be overridden by using the flags.

The main syntax is the following
```
chhoto <subcommand> [<inputs>] [<--flags>]
```

Subcommands:
1.  `new <longurl> [<shorturl>] [<expiry-delay>]`  
    Creates a new shorturl.  
    If shorturl is not provided, it will be generated automatically.  
    Expiry delay should be in seconds. Default value is 0, which means no expiry.
1.  `delete <shorturl>`  
    Deletes a given shorturl.
1.  `expand <shorturl>`  
    Prints info about a particular shorturl.
1.  `getall`  
    Prints info about all shorturls in the server.
1.  `getconfig`  
    Prints the backend config.

Flags:  
    `--api-key`    API Key of the Chhoto URL server.  
    `--password`   Password for the Chhoto URL server. It may also be passed interactively.
    `--url`        URL of the Chhoto URL server.  
    `--version`    Prints the version.
    `-h|--help`    Prints help.

# Notes
- This was just a learning project for me. I wanted to write something tangible using Golang. So, don't expect this to be maintained in the future.
- I haven't used any external packages, everything is written in pure Go, using the Go Standard Library. I'll try to keep it this way.
