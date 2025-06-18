# Chhoto URL CLI
This is a companion CLI tool for [Chhoto URL](https://github.com/SinTan1729/chhoto-url). It provides simple commands for interacting with
a Chhoto URL server.

All the functions of this tool can be replicated using `curl`, but this is supposed to be easier to use. The server needs to have an API key
in order to be accessible using this CLI tool.

# Usage
```
By default, config will be loaded from $XDG_CONFIG_HOME/chhoto/config.json
But these can be overridden by using the flags.
Subcommands:
    new <longurl> [<shorturl>] [<expiry-delay>]    Create a new shorturl.
                                                   If shorturl is not provided, it will be generated automatically.
                                                   Expiry delay should be in seconds. Default value is 0, which means no expiry.
    delete <shorturl>                              Delete a given shorturl.
    expand <shorturl>                              Get info about a particular shorturl.
    getall                                         Get info about all shorturls in the server.
Flags:
    --api-key    API Key of the Chhoto URL server.
    --url        URL of the Chhoto URL server.
    --version    Prints the version.
```
# Notes
- This was just a learning project for me. I wanted to write something tangible using Golang. So, don't expect this to be maintained in the future.
- I haven't used any external packages, everything is written in pure Go, using the Go Standard Library. I'll try to keep it this way.
