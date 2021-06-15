testground run composition -f 30-normal-test.toml --wait
testground run composition -f 30-normal-test.toml --wait
testground run composition -f 30-normal-test.toml --wait
testground run composition -f 30-normal-test.toml --wait
testground run composition -f 30-normal-test.toml --wait
docker system prune --volumes --force
sleep 5
testground run composition -f 30-subburst-test.toml --wait
testground run composition -f 30-subburst-test.toml --wait
testground run composition -f 30-subburst-test.toml --wait
testground run composition -f 30-subburst-test.toml --wait
testground run composition -f 30-subburst-test.toml --wait
docker system prune --volumes --force
sleep 5
testground run composition -f 30-eventburst-test.toml --wait
testground run composition -f 30-eventburst-test.toml --wait
testground run composition -f 30-eventburst-test.toml --wait
testground run composition -f 30-eventburst-test.toml --wait
testground run composition -f 30-eventburst-test.toml --wait
docker system prune --volumes --force
sleep 5
testground run composition -f 30-fault-test.toml --wait
testground run composition -f 30-fault-test.toml --wait
testground run composition -f 30-fault-test.toml --wait
testground run composition -f 30-fault-test.toml --wait
testground run composition -f 30-fault-test.toml --wait
docker system prune --volumes --force
sleep 5
testground run composition -f 30-longrun-test.toml --wait
testground run composition -f 30-longrun-test.toml --wait
testground run composition -f 30-longrun-test.toml --wait
testground run composition -f 30-longrun-test.toml --wait
testground run composition -f 30-longrun-test.toml --wait
docker system prune --volumes --force