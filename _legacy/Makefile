.PHONY: help
.PHONY: conf
.PHONY: deploy
.PHONY: logs

# -- configuration -------------------------------------------------------------
server=h3.tarantsov.com

# ------------------------------------------------------------------------------
help:
	@echo "make reconf            reconfigure server"
	@echo "make deploy            build & deploy to the production server"


# -- server deployment ----------------------------------------------------------

reconf: _deploy/secrets.sh
	cat _deploy/{config.sh,secrets.sh,common.sh,reconf.sh} | ssh $(server) 'bash -s'

deploy: _deploy/secrets.sh
	hugo -D -F
	rsync -avz --delete public/ "$(server):~/tarantsov-draft/"
	cat _deploy/{config.sh,secrets.sh,common.sh,deploy.sh} | ssh $(server) 'bash -s'

_deploy/secrets.sh:
	echo "Copy _deploy/secrets.sh.example into _deploy/secrets.sh and customize!"
	exit 1

logs:
	ssh $(server) 'sudo journalctl -f -u caddy'
