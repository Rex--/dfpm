DFP Manager
=====================
_CLI tool to manage Microchip's Device Family Packs._

This tool downloads CMSIS-Packs from [Microchip Packs Repository](https://packs.download.microchip.com/)
and extracts them to a central install directory located (default `~/.mchp_packs`). It also implements a
search function that allows you to find which pack contains a specific device.

A wrapper around `xc8-cc` allows you to instatiate with just the `-mcpu` flag, automatically searching
for the correct device pack and injecting the `-mdfp` flag into the command.

```
Usage:
  dfpm [-p <path>] install <pack>
  dfpm [-p <path>] list
  dfpm [-p <path>] search [-d] <string>
  dfpm [-p <path>] xc8 [-mcpu=<cpu_id>] [<argument>...]

Options:
  -p <path>, --pack-dir <path>  specify pack directory to use [default: ~/.mchp_packs]
  -d, --download                automatically download packs
```


Examples
---

**Example 1:** Installing device pack for PIC16F1xxxx devices
```
$ dfpm install Microchip.PIC16F1xxxx_DFP.1.31.465.atpack
```

**Example 2:** Searching for and installing device pack containing the pic16lf19197 device
```
$ dfpm search -d 16lf19197
```


Copyright (c) 2026 Rex McKinnon
