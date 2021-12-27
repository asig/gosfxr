# gosfxr

This is a rewrite of [DrPetter's sfxr](https://www.drpetter.se/project_sfxr.html) 
in Go and GTK3 that only exists because I wanted to get my feet wet with UI development
in Go. 
 
Please refer to DrPetter's readme for instructions on how to use it.

## How to build

In order to build `gosfxr`, you need Go 1.13, GTK3, make and Inkscape (`gosfxr` uses 
Inkscape to convert some of the SVG icons to PNGs). 

If you're on Linux, just make sure that you have all the prerequisites installed, 
and then just run

```bash
make all
```

to compile the `gosfxr` binary. Except for GTK3, `gosfxr` has no external dependencies, all
the required resources are statically linked.

If you're on Windows or Mac, the Makefile *might* just work, but I never tested it, and
you're pretty much on uncharted territory :-)

## License
Copyright (c) 2021 Andreas Signer.  
Licensed under [GPLv3](https://www.gnu.org/licenses/gpl-3.0).
