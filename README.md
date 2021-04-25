# coingecko-cache

```                                                              
+--------------------+                                           
|discord-stock-ticker|-----\                                     
|                    |--\   ---------\                           
+--------------------+   -----\       ---------\
                               --> +-----+      ----> +---------+
                                   |redis|            |coingecko|
                               --> +-----+      ----> +---------+
    +----------------+   -----/       ---------/                 
    |coingecko-cache |--/   ---------/                           
    |                |-----/                                     
    +----------------+                                           
```

a cache system for avoiding coingecko rate limits for crypto prices. uses the `coins/markets` endpoint documented [here](https://www.coingecko.com/en/api): 

items are stored in redis under the keys `<coin id>#<attribute>`. requires you have a redis server running locally on port `6379`.

stores info on the top 100 coins by market cap.

## build

```
make build-<linux/osx>
```

## run

```
make run
```

## install

```
make install
```
