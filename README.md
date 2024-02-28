# DwCA is an app and a Go library to deal with Darwin Core Archive files.

Fast reader and writer of Darwin Core Archive Files. For now only
checklist files are supported.

<!-- vim-markdown-toc GFM -->

* [Installation](#installation)
    * [Homebrew on Mac OS X, Linux, and Linux on Windows (WSL2)](#homebrew-on-mac-os-x-linux-and-linux-on-windows-wsl2)
    * [Manual Install](#manual-install)
        * [Linux and Mac without Homebrew](#linux-and-mac-without-homebrew)
        * [Go](#go)
* [Configuration](#configuration)
* [Usage](#usage)
    * [Usage as a command line app](#usage-as-a-command-line-app)
* [Development](#development)
* [Testing](#testing)

<!-- vim-markdown-toc -->

## Installation

### Homebrew on Mac OS X, Linux, and Linux on Windows ([WSL2][wsl])

TLDR:

    ```bash
    brew tap gnames/gn
    brew install dwca
    ```
[Homebrew] is a popular package manager for Open Source software originally
developed for Mac OS X. Now it is also available on Linux, and can easily
be used on MS Windows 10 or 11, if Windows Subsystem for Linux (WSL) is
[installed][wsl].

Note that [Homebrew] requires some other programs to be installed, like Curl,
Git, a compiler (GCC compiler on Linux, Xcode on Mac). If it is too much,
go to the `Linux and Mac without Homebrew` section.

1. Install Homebrew according to their [instructions][Homebrew].

2. Install `dwca` with:

    ```bash
    brew tap gnames/gn
    brew install dwca
    # to upgrade
    brew upgrade dwca
    ```
### Manual Install

`dwca` consists of just one executable file, so it is pretty easy to
install it by hand. To do that download the binary executable for your
operating system from the [latest release][releases].

#### Linux and Mac without Homebrew

Move ``dwca`` executable somewhere in your PATH
(for example ``/usr/local/bin``)

```bash
sudo mv path_to/gnfinder /usr/local/bin
```

#### Go

Install Go v1.22 or higher.

```bash
git clone git@github.com:/gnames/dwca
cd dwca
make tools
make install
```

## Configuration

When you run ``dwca -V`` command for the first time, it will create a
[``dwca.yml``][dwca.yml] configuration file.

This file should be located in the following places:

MS Windows: `C:\Users\AppData\Roaming\dwca.yml`

Mac OS: `$HOME/.config/dwca.yml`

Linux: `$HOME/.config/dwca.yml`

This file allows to set options that will modify behaviour of ``dwca``
according to your needs. It will spare you from entering the same flags for the
command line application again and again.

Command line flags will override the settings in the configuration file.

It is also possible to setup environment variables. They will override the
settings in both the configuration file and from the flags.

| Settings                 | Environment variables           |
|--------------------------|---------------------------------|
| RootPath                 | DWCA_ROOT_PATH                  |
| OutputArchiveCompression | DWCA_OUTPUT_ARCHIVE_COMPRESSION |
| OutputCSVType            | DWCA_OUTPUT_CSV_TYPE            |
| JobsNum                  | DWCA_JOBS_NUM                   |

## Usage

### Usage as a command line app

To see flags and usage:

```bash
dwca --help
# or just
dwca
```

To see the version of its binary:

```bash
dwca -V
```

Normalizing DwCA file

```bash
dwca normalize input_file.zip  <output.zip>
## change number of concurrent jobs
dwca normalize -j 100 input_file.zip  <output.zip>
## change to comma-separated format for the output
dwca normalize -c csv input_dwca.zip
## change to a `tar.gz` archive
dwca normalize -a tar input_dwca.zip
```

If output path is not given, the output will be `{input file name}.norm.zip` or
`{input file name}.norm.tar.gz`

## Development

To install the latest `dwca`

```bash
git clone git@github.com:/gnames/dwca
cd gnfinder
make tools
make install
```

## Testing

To avoid conflicts in filesystem run tests in sequential order.

```bash
go test -p 1 -count=1 ./...
```

[Homebrew]: https://brew.sh/
[wsl]: https://docs.microsoft.com/en-us/windows/wsl/
[dwca.yaml]: https://github.com/gnames/dwca/blob/master/gnfinder/cmd/dwca.yaml
[releases]: https://github.com/gnames/dwca/releases/latest
