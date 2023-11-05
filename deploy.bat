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
