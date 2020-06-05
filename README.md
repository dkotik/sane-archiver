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
sane-archiver keygen
sane-archiver pack [FILE|DIRECTORY]... --key [PUBLICKEY]
sane-archiver unpack [FILE.sane1]... --key [PRIVATEKEY]
sane-archiver --help [keygen|pack|unpack]
```

    Options:
     -o, --output       Output to this file or path.
     -f, --force        Overwrite any files that already exist.
     -w, --warn <GB>    Warn if the disk is running low on space.
     -l, --leave <X>    Delete all *.sane1 files except X most
                        modified most recently.

     Defaults:
       --output defaults to {year}-{month}-{day}-{md5}.[sane1|zip]
       --key [PUBLICKEY] defaults to $ENV[SaneArchiverPublicKey]
       --warn defaults to 2, issuing a warning under 2GB of free space

<!-- The path to newly created archive is printed into os.Stdout. The log of the creation process and any warnings or errors are printed into os.Stderr. This simplifies the creation of recipes that log to a certain file or notify you by email when archives are created or upload created files:

```bash
sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY] 2>>report.log
sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY] 2>>&1 | tee report.log | mail -s "Email subject" me@mymail.com
sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY] 2>>report.log && aws s3 [DIRECTORY] sync s3://bucket...
``` -->

**Pro tip:** If you type one space character before running a terminal command, that command will not be recorded in your bash history. This can help protect your keys from prying eyes.

Features
--------

*   **No Artifacts**. Archiver produces files that appear to contain entirely
    random-generated data without any markings or artifacts. Stream cipher is used to
    encrypt the containing data. If a produced archive file
    ends up in the hands of a malicious actor, it will be difficult to determine
    how the file was created just by looking at its contents or its size. Do not forget to change
    the default file-naming scheme by using `--output {hash}.extension` command line argument.

*   **Git Archive Support**. Archiver detects folders that contain Git repositories and archives
    all Git branches as separate *.tar balls. (Requires Git to be installed on the machine!)

*   **S3 Upload**. Archiver can attempt to upload the resulting file to AWS S3 upon completion.
    Use `--upload s3://<credentialID>:<credentialSecret>@<awsRegion>/<bucket>/<path>` parameter.
    Note that the local copy of the file will be retained. You can protect your disk from filling up by accident by setting `--output /tmp/{hash}.tmp`.

*   **Includes MD5 Hash In Output**. By default, generated files include MD5 hash in their name.
    Thus, checking for bit-rot errors is as trivial as running `md5sum .`

*   **File System Warnings**. Archiver will print a warning if the target file system
    is running low on available storage space. By default, the warning is printed when there
    are less than 2GB of space remains. You can change the warning threshold by passing
    `--warn [INTEGER]` as a command line argument.

Roadmap
-------

- Check if s3://URL is a directory, UploadS3 does not work if URL points to a directory.
- Checksum --md5 command.
- Support git sub-modules for archiving. Currently they are ignored.
- [Stash git changes](https://stackoverflow.com/questions/2766600/git-archive-of-repository-with-uncommitted-changes) before making an archive.
- Add support for Windows (Linux and MacOS are both supported).
- Display progress percentage when running through files and when uploading.

License
-------

Sane Archiver is distributed under Apache License. The author would also like to add the SQLite blessing:

> May you do good and not evil. May you find forgiveness for yourself and forgive others. May you share freely, never taking more than you give.
