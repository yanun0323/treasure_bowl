# Treasure Bowl

A simple crypto trading bot.

# Requirement
#### _Go 1.21.4 or higher_

# Structure Design
### General structure
```mermaid
graph TD
    KlineProvider{"K-Line Provider\nInterface"}
    TradeServer{"Trade Server\nInterface"}
    TradeServer2{"Trade Server\nInterface"}
    STG(("Strategy Logic"))

    KlineProvider ==>|"get & listening price changing"| STG 
    TradeServer ==>|"get asset & order in the beginning\nlistening asset/order changing"| STG

    STG ==>|"push order according to strategy\nget changed account asset"| TradeServer2

    TradeServer -.-|same| TradeServer2
```

# TODO
- KlineProvider: victor 研究中
- AssetProvider
- Strategy
- Order