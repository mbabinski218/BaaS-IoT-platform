DATADIR=/data

if [ ! -d "$DATADIR/geth" ]; then
  echo ">>> Inicjalizacja genesis.json..."
  geth init /genesis/genesis.json --datadir $DATADIR
fi

echo ">>> Startujemy geth..."
exec geth \
  --datadir $DATADIR \
  --networkid 1234 \
  --mine \
  --http \
  --http.addr 0.0.0.0 \
  --http.api web3,eth,personal \
  --http.port 8545 \
  --http.vhosts "*" \
  --nodiscover \
  --allow-insecure-unlock