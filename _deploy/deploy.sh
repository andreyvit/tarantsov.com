#!/bin/bash
set -x -euo pipefail

username="andreyvit"
now="$(date "+%Y%m%d_%H%M%S")"

# -- deploy --------------------------------------------------------------------

echo "▸ deployment"
sudo install -d -m755 -g$username -o$username /srv/tarantsov-www/versions
sudo cp -r ~/tarantsov-www "/srv/tarantsov-www/versions/$now"
sudo chown -R $username:$username "/srv/tarantsov-www/versions/$now"
sudo rm -f "/srv/tarantsov-www/upcoming"
sudo ln -s "/srv/tarantsov-www/versions/$now" "/srv/tarantsov-www/upcoming"
sudo mv -T "/srv/tarantsov-www/upcoming" "/srv/tarantsov-www/current"

# -- Caddyfile -----------------------------------------------------------------
echo "▸ Caddyfile"
sudo install -m644 -groot -oroot /dev/stdin /srv/tarantsov-www/Caddyfile <<EOF
tarantsov.com {
    root /srv/tarantsov-www/current
    tls andrey@tarantsov.com
}
EOF

echo "✓ all done"
