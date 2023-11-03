# api

test token
```txt
v4.public.eyJleHAiOiIyMDIzLTEyLTAzVDE3OjQyOjQ0KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0wM1QxNzo0Mjo0NCswNzowMCIsImlkIjoiYXdhbmdnYSIsIm5iZiI6IjIwMjMtMTEtMDNUMTc6NDI6NDQrMDc6MDAifWMmQMzPu5CpQBm3Id0-ODEoWvwK_ABO2cw_07OXzpUtp1j2cdXEyimRvzCunJjcDCCqLbFNcqSZZxelcOzpKww
```

## build

```sh
$env:GOOS = 'linux'
$env:CGO_ENABLED = '1'
go build
scp -P 24520 api root@103.155.250.23:/root/api
ssh -p 24520 root@103.155.250.23 chmod +x /root/api
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