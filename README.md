## Hardware based license server

forked and modified from https://github.com/SaturnsVoid/HWID-Based-License-System

## Usage

### Buld and run server

```bash
docker network create lic
docker build -t licserver . 
docker run \
    --network lic \
    -t -d --rm -p 9347:9347 \
    -v $(pwd)/db:/app/db \
    --name licserver  \
    -e SALT="12345salt"  \
    -e TOKEN="mytoken" \
    licserver
```

### Add first license
```bash
docker exec -it licserver /app/server
```

### Run client
```bash
License=86UU-N4SB-OQYH go run client/client.go 
```

## Try to run client in Docker
```bash
docker build -t client -f Dockerfile-client . 
docker run --network lic -t --name client -e License="16FB-L6AX-2ZPZ" client
```

```bash
docker start client
docker logs client | tail -n 2
```



## API
### ADD
```bash
curl -X POST http://127.0.0.1:9347/add \
-H "Authorization: mytoken" \
-H "Content-Type: application/json" \
-d '{"email": "test1", "expiration": "2023-12-31"}'

{"email":"test1","exp_date":"2023-12-31","license":"16FB-L6AX-2ZPZ","message":"New license generated"}
```

### LIST
```bash
curl -X GET http://127.0.0.1:9347/list \
-H "Authorization: mytoken"

["16FB-L6AX-2ZPZ:2023-12-31:test1:NOTSET"]
```

### RESET
```bash
curl -X POST http://127.0.0.1:9347/reset-key \
-H "Authorization: mytoken" \
-H "Content-Type: application/json" \
-d '{"key": "16FB-L6AX-2ZPZ"}'
    
    {"email":"test1","exp_date":"2023-12-31","license":"16FB-L6AX-2ZPZ","message":"License is valid"}
```

### REMOVE
```bash
curl -X POST http://127.0.0.1:9347/remove
-H "Authorization: mytoken" \
-H "Content-Type: application/json" \
-d '{"email": "test1"}'

{"email":"test1","message":"License removed"}
```


{"email":"test1","exp_date":"2023-12-31","license":"16FB-L6AX-2ZPZ","message":"License is valid"}
```



_________________________________________________________________________________________

# HWID-Based-License-System
A GoLANG based HWID license system, basic.

Vary simple, basic HWID (hardware ID) license system.

You generate keys with the license server and give the program to the client with a key, on first run the program looks for the license.dat file that contains the key if not found asks the client if they want to register, if so they imput the key and the program generates a HWID for that system and user, connects to the license server where the server checks for the key making sure its not already registerd with another HWID and that its not expired. If good it adds the HWID to the row in the database.

You can generate new keys will the following information;

Email
Experation Date

The key is generated from a random char generater set to 4x4x4 chars 0-9 and A-Z.

You can also bulk generate keys (without registerd email)
You also can remove keys (by email)

This is not fail proof, its a simple to use deterent.

The database is just a text file, ity can be edited by hand.

THE PROGRAM WILL NEED TO BE ABLE TO CONNECT TO THE SERVER TO VERIFY THE LICENSE ANYTIME ITS CALLED.
THE LICENSE CHECK CAN BE RUN AT ANYTIME OR ON A TIMED LOOP.
LICENSE SERVER CAN RUN ON ANY OPEN PORT.

# Donations
<img src="https://blockchain.info/Resources/buttons/donate_64.png"/>
<p align="center">Please Donate To Bitcoin Address: <b>1AEbR1utjaYu3SGtBKZCLJMRR5RS7Bp7eE</b></p>
