DATADIR=/data

echo ">>> Begin geth..."

if [ ! -d "$DATADIR/geth" ]; then
  echo ">>> Init geth..."
  cat /genesis/genesis.json
  geth --datadir $DATADIR init /genesis/genesis.json
else
  echo ">>> Skipping init..."
fi

echo ">>> Start geth..."
exec geth \
  --datadir $DATADIR \
  --networkid 1234 \
  --mine \
  --http \
  --http.addr 0.0.0.0 \
  --http.api web3,eth,net  \
  --http.port 8545 \
  --http.vhosts "*" \
  --nodiscover \
  --syncmode "snap" \
  --cache 4096 \
  --maxpeers 0 \
  --allow-insecure-unlock