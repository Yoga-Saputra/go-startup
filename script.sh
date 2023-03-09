git pull origin main
rm -rf /var/www/go-startup/startup
export PATH=$PATH:/usr/local/go/bin
go build .
sudo systemctl restart startup
