# api

## disable not used services

```sh
systemctl disable exim4
systemctl disable apache2
systemctl disable xinetd
```

## build

```sh
$env:GOOS = 'linux'
go build
```

## Create Service

```sh
sudo vim /lib/systemd/system/api.service
```

```conf
[Unit]
Description=API Server

[Service]
Type=simple
Restart=always
RestartSec=10s
ExecStart=/root/api

[Install]
WantedBy=multi-user.target
```

```sh
sudo vim /lib/systemd/system/api.service
```

```sh
systemctl enable fiberserver
systemctl start fiberserver
reboot
```

## ENV

```sh
vim /etc/environment
```
