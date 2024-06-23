# Secret Service

As of **kwalletd6** is officially released, this project is obsolete.

[![GitHub release](https://img.shields.io/github/release/yousefvand/secret-service.svg?style=plastic)](https://github.com/yousefvand/secret-service/releases)
[![GitHub license](https://img.shields.io/github/license/yousefvand/secret-service.svg?style=plastic)](https://github.com/yousefvand/secret-service/blob/master/LICENSE.md)
[![GitHub stars](https://img.shields.io/github/stars/yousefvand/secret-service.svg?style=plastic)](https://github.com/yousefvand/secret-service/stargazers)
[![GitHub issues](https://img.shields.io/github/forks/yousefvand/secret-service.svg?style=plastic)](https://github.com/yousefvand/secret-service/forks)
[![GitHub issues](https://img.shields.io/github/issues/yousefvand/secret-service.svg?style=plastic)](https://github.com/yousefvand/secret-service/issues)

Implementation of [Secret Service API](http://standards.freedesktop.org/secret-service)

![logo](assets/secret-service.png)

## What does this project do?

By using **secret service**, you don't need to use `KeePassXC` _secretservice_ for storing and retrieving you applications credentials anymore, or login every time to `Skype`, `vscode sync`, `Remmina`...

## Installation

- Archlinux: There is an [AUR package](https://aur.archlinux.org/packages/secret-service/) named `secret-service`.
- Debian: _TODO_ deb package
- RedHat: _TODO_ rpm package

## Manual Installation

There is a `scripts/manage.sh` shellscript that do the job of install/uninstall (run it by `./scripts/manage.sh`) but here are the details:

You need to copy the binaries (`secretserviced` and `secretservice`, build the project or download it from [releases](https://github.com/yousefvand/secret-service/releases) page) some where usually `/usr/bin` but if you don't have the permission, `~/.local/bin` is OK too. To build the binaries from source code:

```bash
git clone https://github.com/yousefvand/secret-service.git
cd secret-service
go build -race -o secretserviced cmd/app/secretserviced/main.go
go build -race -o secretservice cmd/app/secretservice/main.go
```

You need a `systemd` **UNIT** file named `secretserviced.service` to put in `/etc/systemd/user` but if you don't have the permission `~/.config/systemd/user` is OK too. Here is a sample **UNIT** file, change `WorkingDirectory` and `ExecStart` according to where you put the binary (`secretserviced`):

```config
[Unit]
Description=Service to keep secrets of applications
Documentation=https://github.com/yousefvand/secret-service

[Install]
WantedBy=default.target

[Service]
Type=simple
RestartSec=30
Restart=always
Environment="MASTERPASSWORD=01234567890123456789012345678912"
WorkingDirectory=/usr/bin/
ExecStart=/usr/bin/secretserviced
```

**CAUTION**: `MASTERPASSWORD` is very important, don't loose it. `scripts/manage.sh` would generate a random `32` character password automatically. If you don't use the `scripts/manage.sh` shellscript, it is up to you to set the password and it should be **EXACTLY** `32` characters length.

Now start the service:

```bash
sudo systemctl daemon-reload
systemctl enable --now --user secretserviced.service
```

and you can stop the service by:

```bash
systemctl disable --now --user secretserviced.service
```

to see the status of service:

```bash
systemctl status --user secretserviced.service
```

All `secret-service` stuff (database, logs...) are stored under: `~/.secret-service`.

By default all secrets are encrypted with `AES-CBC-256` symmetric algorithm with `MASTERPASSWORD`. If you wish to switch between encrypted/unencrypted database you need to follow these steps:

1. Stop service: `systemctl stop --user secretserviced.service`
2. Change config `encryption` key (located at: `~/.secret-service/secretserviced/config.yaml`)
3. If you are changing to `encryption: true` make sure `MASTERPASSWORD` is set.
4. Delete database (located at: `~/.secret-service/secretserviced/db.json`)
5. Start service: `systemctl start --user secretserviced.service`

If service refuses to start and you see `OS` exit code `5` in logs, it means som other application has taken dbus name `org.freedesktop.secrets` before (such as keyrings), stop that application and try again.

## secretservice

This binary is the `CLI` interface to communicate with `secretserviced` daemon. Supported commands:

### ping

```bash
secretservice ping
```

Check if service is up and responsive.

### export db

```bash
secretservice export db
```

Export a copy of current db in `~/.secret-service/secretserviced/`. This copy is not encrypted.

### encrypt

```bash
secretservice encrypt -p|--password 32character-password -i|--input /path/to/input/file/ -o|--output /path/to/output/file/
```

Encrypts input file using given password. Password should be exactly 32 character. Example:

```bash
secretservice encrypt -p 012345678901234567890123456789ab -i ~/a.json -o ~/b.json
```

### decrypt

```bash
secretservice decrypt -p|--password 32character-password -i|--input /path/to/input/file/ -o|--output /path/to/output/file/
```

Decrypts input file using given password. Password should be exactly 32 character. Example:

```bash
secretservice decrypt -p 012345678901234567890123456789ab -i ~/a.json -o ~/b.json
```

## Contribution

This project is in its infancy and as it is my first golang project there are many design and code problems. I do appreciate suggestions and **PR**s. If you can get done any item from `TODO` list, you are welcome. This list will be updated based on new insights and user issues.

In case of sending a **PR** please make sure:

1. You are addressing just one issue per PR.
2. Completely describe the problem and your solution in plain English.
3. Don't send your PRs to `main` branch, create a new branch based on your changes and make sure all tests are passed.
4. If any new test is needed based on your PR, please write the test as well.

### TODO

- [ ] Improve CI

- [ ] What's the best way to secure `/etc/systemd/user/secretserviced.service` file

- [ ] deb, rpm, AppImage packages

- [ ] ...
