# Word of wisdom

## How to build and run

Clone the source code
```
git clone https://github.com/maksadbek/workofwisdom.git
```

Run with docker-compose

```
docker-compose up -d server
```

you must see that docker-compose is running the server

Run client

```
docker-compose rm --rm client sh
```

This command opens the shell. In the shell, run client:

```
/app/client
```

client will generate valid hashcash payload, and send it to server, then write the response to stdout:

```
2024/05/12 16:54:06 generated a token in 1.131505375s: "X-Hashcash: 1:20:240512165405:username::c1JDR1d4UE1Lcw==:5OtZ"
2024/05/12 16:54:06 received message: "In the midst of chaos, there is also opportunity. - Sun Tzu"
```
