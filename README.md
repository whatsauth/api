# api

test token

```txt
v4.public.eyJleHAiOiIyMDIzLTEyLTAzVDE4OjEwOjI0KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0wM1QxODoxMDoyNCswNzowMCIsImlkIjoiNjI4Nzc1MjAwMDMwMCIsIm5iZiI6IjIwMjMtMTEtMDNUMTg6MTA6MjQrMDc6MDAifY0hFvyfUYFyAIHNZe-XolWxtU5ChlhxJEwXMdqnHMBIexQ00Ew88XGWmltFZGdo0m_ekhpV2oJElGR6eTK80Aw
```

```txt
v4.public.eyJleHAiOiIyMDIzLTEyLTA0VDA5OjE1OjE0KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0wNFQwOToxNToxNCswNzowMCIsImlkIjoiNjI4MzEzMTg5NTAwMCIsIm5iZiI6IjIwMjMtMTEtMDRUMDk6MTU6MTQrMDc6MDAifSqR5kBfQhwRfrtrMiOxXNoPP0syIUPpEbtOMqdPOMEfXbOC6boO6NDFKCKKSqjY8WfTcDBXAHtC9N7NHjrvmwM
```

## build

CGO_ENABLED=0
go build -a

```sh
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api .
$env:GOOS = 'linux'
$env:CGO_ENABLED = '1'
go build
scp -P 24520 api root@103.155.250.23:/root/apinew
ssh -p 24520 root@103.155.250.23 chmod +x /root/apinew
ssh -p 24520 root@103.155.250.23 ls -l
ssh -p 24520 root@103.155.250.23 systemctl status api
ssh -p 24520 root@103.155.250.23 systemctl stop api
ssh -p 24520 root@103.155.250.23 systemctl status api
ssh -p 24520 root@103.155.250.23 mv /root/apinew /root/api
ssh -p 24520 root@103.155.250.23 systemctl start api
ssh -p 24520 root@103.155.250.23 systemctl status api

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

## NGINX Conf

```conf
server {
    listen 443 ssl;

    server_name api.wa.my.id;

    ssl_certificate /etc/letsencrypt/live/api.wa.my.id/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.wa.my.id/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;
 
 client_max_body_size 999M;
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header X-Real-IP  $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
        proxy_set_header X-Forwarded-Port 443;
        proxy_set_header Host $host;
    }
}
```

## Postgresql

```sh
/root/.fly/bin/flyctl auth login
/root/.fly/bin/fly postgres create
/root/.fly/bin/flyctl postgres connect -a whatsauth
/root/.fly/bin/flyctl proxy 5432 -a whatsauth 
/root/.fly/bin/flyctl postgres restart -a whatsauth
```

create db.sh

```sh
#!/bin/sh
while :
do
  /root/.fly/bin/flyctl proxy 5432 -a whatsauth
done
```

```sh
crontab -e
```

```conf
@reboot /root/db.sh
```
