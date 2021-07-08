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
  -db int
        redis db to use
  -expiry int
        number of seconds to keep entries in the cache (default 60)
  -frequency int
        seconds between updates (default 1)
  -hostname string
        connection address for redis (default "localhost:6379")
  -pages int
        number of pages (100 coin each) to pull from (default 1)
  -password string
        redis password
```

```
make run
```

## install

```
make install
```
