Sane Archiver Alpha
===================

[![Go Report Card](https://goreportcard.com/badge/github.com/dkotik/sane-archiver)](https://goreportcard.com/report/github.com/dkotik/sane-archiver)

Sane Archiver is a simple command line utility for making encrypted archives. Rsync is [still one of the most common tools](https://www.tecmint.com/linux-system-backup-tools/) for making and moving backups, which indicates the need for a simpler tool that builds encrypted archives and integrates well with other command line tools. Sane Archive follows Linux philosophy and draws inspiration from [WireGuard](https://www.wireguard.com/) by making user choices limited. Future releases of Sane Archiver will not be backwards-compatible, although the ability to decrypt archives made by previous versions will be maintained, mostly. This program is distributed under the [Apache License Version 2.0](LICENSE). This is the author's first free open source project.


Installation
------------

With Go and git installed:

```bash
go version  # verify that Go Lang is installed
git clone https://github.com/dkotik/sane-archiver.git
cd sane-archiver
sudo make linux install
```

Usage
-----

```bash
sane-archiver --keygen
sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY]... [OPTION]...
sane-archiver --key [PRIVATEKEY] --decrypt [SANEFILE]
```

    Options:
     -k, --key <KEY>    Set private or public base64-encoded key.
         --keygen       Generate a base64-encoded keypair.
     -o, --output       Output to this file or path.
     -d, --decrypt      Decrypt this file using the key.
     -f, --force        Overwrite any files that already exist.
     -w, --warn <GB>    Warn if the disk is running low on space.
     -h, --help         Print this message.

     Defaults:
       --output defaults to {year}-{month}-{day}-{md5}.[sane1|zip]
       --key [PUBLICKEY] defaults to $ENV[SaneArchiverPublicKey]
       --warn defaults to 2, issuing a warning under 2GB of free space

The path to newly created archive is printed into os.Stdout. The log of the creation proccess and any warnings or errors are printed into os.Stderr. This simplifies the creation of recipes that log to a certain file or notify you by email when archives are created or upload created files:

```bash
sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY] 2>>report.log
sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY] 2>>&1 | tee report.log | mail -s "Email subject" me@mymail.com
sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY] 2>>report.log && aws s3 [DIRECTORY] sync s3://bucket...
```

Features
--------

*   **No Artifacts**. Sane Archiver produces files that appear to contain entirely
    random-generated data without any markings or artifacts. Stream cipher is used to
    encrypt the containing data. If a produced archive file
    ends up in the hands of a malicious actor, it will be difficult to determine
    how the file was created just by looking at its contents or its size. Do not forget to change
    the default file-naming scheme by using `--output {hash}.extension` command line argument.

*   **Includes MD5 Hash In Output**. By default, generated files include MD5 hash in their name.
    Thus, checking for bit-rot errors is as trivial as running `md5sum .`

*   **File System Warnings**. Sane Archiver will print a warning if the target file system
    is running low on available storage space. By default, the warning is printed when there
    are less than 2GB of space remains. You can change the warning threshold by passing
    `--warn [INTEGER]` as a command line argument.

TODO
----

- Add support for Windows and MacOS.
- Provide test data and write a beefier test suite.
- Display progress percentage when running through files.
