
# Incentivized Liquidation System for Perpetual Market

## Diagram

```plaintext
+------------+                  +-------------+                  
|            |                  |             |                
| Borrowers  |                  | Watchers    |                
|            |                  |             |                  
+-----+------+                  +-----+-------+               
      |                               |                             
      | Borrow funds to open          | Check positions                
      | positions                     | for liquidation or stop loss conditions     
      |                               |                              
+-----v------+                  +-----v-------+                   
|            |                  |             |                 
| Perpetual  |<-----------------| Liquidation/|
| Position   |   Monitor        |    Close    |   Monitor        
| Module     |   positions      |    Function |   positions       
+------------+                  +-------------+                  
      |                               |                               
      |                               |                                
      | MTP health calculated         | 1. If health < threshold, trigger liquidation
      |                               | 2. If stop loss conditions are met, close position       
      |                               | and receive reward             
      |                               |                                  
+-----v------+                  +-----v-------+                  
|            |                  |             |                   
|  Borrower  |                  |  Watchers   |                   
| Manages    |                  | Earns       |                  
| Position   |                  | Reward      |                  
+------------+                  +-------------+                   
         
```

## Explanation

- **Borrowers** take perpetual positions by borrowing from the pool.
- **Watchers**: Watchers are users and bots that monitor positions for liquidation or stop loss conditions.
- **Bots**: Bots are automated programs specifically designed to efficiently handle these tasks.
- The system rewards liquidators and bots for maintaining the health of the market.
- The **Liquidation Function** is responsible for managing the logic of liquidations and distributing rewards.

This approach distributes the load of monitoring and liquidation to incentivized users and bots, making the system more efficient and scalable.
