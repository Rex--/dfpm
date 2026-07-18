DFP Manager
=====================
_CLI tool to manage Microchip's Device Family Packs._

This tool downloads CMSIS-Packs from [Microchip Packs Repository](https://packs.download.microchip.com/)
and extracts them to a central install directory located at `~/.mchp_packs`. It also implements a
search function that allows you to find which pack contains a specific device.

A wrapper around `xc8-cc` allows you to instatiate with just the `-mcpu` flag, automatically searching
for the correct device pack and injecting the `-mdfp` flag into the command.

```
Usage: dfpm [install | list | search ] [arg]
       dfpm install <pack>      install the given device pack
       dfpm list                list installed device packs
       dfpm search <string>     search for <string> in pack descriptions

       dfpm xc8 [-mcpu=<cpuid>] [args ...]
                                search for the given <cpuid> in pack descriptions
                                if found, launch 'xc8-cc' with the -mdfp option
```

Copyright (c) 2026 Rex McKinnon
