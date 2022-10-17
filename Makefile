server=rhodes.tarantsov.com

ship:
	hugo --minify
	rsync -avz --delete public/ "$(server):~/tarantsov-www/"
	ssh $(server) bash -s -- <_deploy/deploy.sh

reload:
	ssh $(server) /srv/caddy/bin/caddy reload --config /etc/Caddyfile

logs:
	ssh $(server) 'sudo journalctl -f -u caddy'
