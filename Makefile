server=h3.tarantsov.com

ship:
	hugo --minify
	rsync -avz --delete public/ "$(server):~/tarantsov-www/"
	ssh $(server) bash -s -- <_deploy/deploy.sh

restart:
	ssh $(server) sudo systemctl restart caddy

logs:
	ssh $(server) 'sudo journalctl -f -u caddy'
