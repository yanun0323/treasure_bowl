# Treasure Bowl

A simple crypto trading bot.

# Requirement
#### _Go 1.21.4 or higher_

# System Design
```mermaid
graph TD
    PriceProvider{"Price Provider\nInterface"}
    AssetProvider{"Asset Provider\nInterface"}
    Strategy{Strategy\nInterface}
    Order{Order\nInterface}

    PriceProvider --->|1. Listen Price Changing| BotService
    AssetProvider -->|2. Listen Asset Changing| BotService
    BotService -->|3. Send Price/Asset| Strategy 
    Strategy -->|4. Return Signal| BotService
    BotService --->|5. Create/Cancel Order| Order
```