# GoFast
![gofast-CI](https://github.com/AandersonL/gofast/workflows/gofast-CI/badge.svg)

A HTTP downloader accelerator for Linux like [axel](https://github.com/axel-download-accelerator/axel), but faster.

## What is this ?

gofast is a http-only downloader that uses multiple connections through concurrency relying in [Range Request](https://tools.ietf.org/rfc/rfc7233.txt). It can burst your download speed just by spliting the whole GET request in ***N*** connections.

## Features

* Normal HTTP download if the server does not accept ranges
* HTTP concurrency download
* Control of how many parallel connections per host

## Usage
There is a couple of options that you can use:

<pre>
$ gofast -h                     
usage: gofast [-h|--help] [-n|--num-tasks <integer>] [-o|--output "<value>"]
              -u|--url "<value>"

              A HTTP downloader accelerator using concurrency

Arguments:

  -h  --help       Print help information
  -n  --num-tasks  Number of concurrent connections, default: Num cores * 2.
                   Default: 0
  -o  --output     Output filename. Default: 
  -u  --url        Remote file url

</pre>

Simple usage:

<pre>
$ gofast -u https://releases.ubuntu.com/20.04.1/ubuntu-20.04.1-desktop-amd64.iso                                                         
Starting concurrent download of ubuntu-20.04.1-desktop-amd64.iso
Downloading 100% |...| (2.6/2.6 GB, 10.634 MB/s)          
Downloaded in 4m9.761532636s
</pre>

Downloading the same file using axel:

<pre>
$ axel -a https://releases.ubuntu.com/20.04.1/ubuntu-20.04.1-desktop-amd64.iso
Initializing download: https://releases.ubuntu.com/20.04.1/ubuntu-20.04.1-desktop-amd64.iso
File size: 2785017856 bytes
Opening output file ubuntu-20.04.1-desktop-amd64.iso.0
Starting download

Connection 3 unexpectedly closed
Connection 0 finished
Connection 1 finished
Connection 2 finished
Connection 1 finished
Connection 3 finished
Connection 0 finished

Downloaded 2.59375 Gigabyte(s) in 4:32 minute(s). (9969.09 KB/s)
</pre>




## Installing

You can go to [releases](https://github.com/AandersonL/gofast/releases) and grab the latest one, or manually build.

> $ git clone https://github.com/aandersonl/gofast && cd gofast && go build


Any bugs or enhancement feel free to open a issue or pull request!
