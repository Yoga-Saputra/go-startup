git pull origin main
rm -rf /var/www/go-startup/startup
export PATH=$PATH:/usr/local/go/bin
go build .
# sudo service startup stop
# sudo service startup start
systemctl daemon-reloads
