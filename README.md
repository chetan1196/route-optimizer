## Overview

The Delivery System assigns delivery executives batches of orders and calculates the best route for delivering them in the shortest possible time. It utilizes two strategies to achieve this: the Concurrent Naive Strategy and the Dynamic Programming Strategy.

## Features

- Calculates optimal delivery routes for batches of orders
- Utilizes the Concurrent Naive Strategy for simple route calculation with concurrency
- Implements the Dynamic Programming Strategy for more efficient route calculation
- Considers travel time between geo-locations using the Haversine formula

## Usage

- `go build`
- `./route-optimizer`

## Sample Output
./route-optimizer 
2024/05/21 12:40:40 Total travel time: 960.82 minutes
2024/05/21 12:40:40 Step 1: Pick up from Restaurant B at (12.9820, 13.6700)
2024/05/21 12:40:40 Step 2: Deliver to Consumer B at (12.9370, 12.8940)
2024/05/21 12:40:40 Step 3: Pick up from Restaurant A at (12.0820, 13.2700)
2024/05/21 12:40:40 Step 4: Deliver to Consumer A at (12.9160, 12.5940)
