DFP Manager
=====================
_CLI tool to manage Microchip's Device Family Packs._

This tool downloads CMSIS-Packs from [Microchip Packs Repository](https://packs.download.microchip.com/)
and extracts them to a central install directory (default `~/.mchp_packs/Microchip`).
It allows you to manage these packs through listing, removing, and updating them.
It also implements a search function that allows you to find which pack contains a specific device.

A wrapper around `xc8-cc` allows you to invoke with just the `-mcpu` flag, automatically searching
for the correct device pack and injecting the `-mdfp` flag into the command.
Ensure that `xc8-cc` is available through the `PATH` enviornment variable.

```
dfpm - Device Family Pack Manager
Download, install, and search Microchip DFPs.

  Usage:
    dfpm [install|remove|list|search|update|xc8]

  Subcommands:
    install   Install a given device family pack.
    remove    Remove a given device family pack.
    list      List device family packs.
    search    Search for a given term in device packs.
    update    Update installed pack(s).
    xc8       Launch xc8-cc with -mdfp argument based on -mcpu.

  Flags:
        --version    Displays the program version string.
    -h  --help       Displays help with available flag, subcommand, and positional value parameters.
    -p  --pack-dir   Specify pack directory to use. (default: ~/.mchp_packs/Microchip)

Call each subcommand with the help flag for more inforation about the usage.
```


Examples
---

**Example 1:** Installing device pack for PIC16F1xxxx devices
```
$ dfpm install pic16f1xxxx
downloading https://packs.download.microchip.com/Microchip.PIC16F1xxxx_DFP.1.32.471.atpack
unpacking Microchip.PIC16F1xxxx_DFP.1.32.471.atpack
installed /home/dfpm/.mchip_packs/Microchip/Microchip.PIC16F1xxxx_DFP.1.32.471
```

**Example 2:** Removing device packs
```
$ dfpm remove pic16f1
Found 1 packs matching 'pic16f1'
  PIC16F1xxxx_DFP v1.32.471
Remove 1 packs? [Y/n]: y
Removing /home/dfpm/.mchip_packs/Microchip/Microchip.PIC16F1xxxx_DFP.1.32.471
```

**Example 3:** Listing installed device packs
```
$ dfpm list
PIC10-12Fxxx_DFP  v1.10.193
PIC12-16Cxxx_DFP  v1.9.187
PIC12-16F1xxx_DFP v1.10.270
PIC16F1xxxx_DFP   v1.32.471
PIC16Fxxx_DFP     v1.9.176
PIC18F-Q_DFP      v1.31.492
```

**Example 4:** Searching for and installing device pack containing the pic16lf19197 device
```
$ dfpm search -d 16lf19197
Found 'pic16lf19197' in 1 packs
  PIC16F1xxxx_DFP
Install 1 packs? [Y/n]: y
downloading https://packs.download.microchip.com/Microchip.PIC16F1xxxx_DFP.1.32.471.atpack
unpacking Microchip.PIC16F1xxxx_DFP.1.32.471.atpack
installed /home/dfpm/.mchip_packs/Microchip/Microchip.PIC16F1xxxx_DFP.1.32.471
```

**Example 5:** Using xc8 to compile a blink program. The device pack will be downloaded if not available.
```
$ dfpm xc8 -mcpu=16lf19197 blink.c -o blink.hex
Found device 16lf19197 in pack PIC16F1xxxx_DFP, installing...
downloading https://packs.download.microchip.com/Microchip.PIC16F1xxxx_DFP.1.32.471.atpack
unpacking Microchip.PIC16F1xxxx_DFP.1.32.471.atpack
installed /home/dfpm/.mchip_packs/Microchip/Microchip.PIC16F1xxxx_DFP.1.32.471

16LF19197 Memory Summary:
    Program space        used    25h (    37) of  8000h words   (  0.1%)
    Data space           used     3h (     3) of  1000h bytes   (  0.1%)
    EEPROM space         used     0h (     0) of   100h bytes   (  0.0%)
    Configuration bits   used     5h (     5) of     5h words   (100.0%)
    ID Location space    used     0h (     0) of     4h bytes   (  0.0%)
```

Copyright © 2026 Rex McKinnon \
This software is released under the permissive University of Illinois/NCSA Open Source License.
Check the LICENSE file for more details.
