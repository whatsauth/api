# api

## build

```sh
$env:GOOS = 'linux'
go build
scp -P 24520 api root@103.155.250.23:/root/temp
ssh -p 24520 root@103.155.250.23 chmod +x /root/temp
ssh -p 24520 root@103.155.250.23 systemctl stop api
ssh -p 24520 root@103.155.250.23 mv /root/temp /root/api
ssh -p 24520 root@103.155.250.23 systemctl start api
```

## disable not used services

```sh
systemctl disable exim4
systemctl disable apache2
systemctl disable xinetd
```

## Create Service

```sh
vim /lib/systemd/system/api.service
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
systemctl enable api
systemctl start api
reboot
```

## ENV

```sh
vim /etc/environment
```
