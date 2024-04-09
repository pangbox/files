# PangFiles [![godoc.org](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/pangbox/pangfiles) [![Go Report Card](https://goreportcard.com/badge/github.com/pangbox/pangfiles)](https://goreportcard.com/report/github.com/pangbox/pangfiles)

PangFiles is a set of Go libraries that implement various PangYa file formats, encryption schemes, and hashes.

# Nix Flake

On UNIX-like systems with Nix installed, you can use the Nix flake to build and run PangFiles.

```shell
$ nix run github:pangbox/pangfiles pak-mount *.pak ~/mnt
```

To install the Nix package manager, use the [Determinate Nix Installer][1].

[1]: https://github.com/DeterminateSystems/nix-installer "The Determinate Nix Installer"
