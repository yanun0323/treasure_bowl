.PHONY:

PAIR=BTC_USDT
STG=inspector
KLINE=bitopro,mock
KLINE_DURATION=1m
TRADE=mock

run:
	PAIR=${PAIR} \
	STG=${STG} \
	KLINE=${KLINE} \
	KLINE_DURATION=${KLINE_DURATION} \
	TRADE=${TRADE} \
	go run ${CURDIR}/main.go