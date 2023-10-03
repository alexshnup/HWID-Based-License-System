## Hardware based license server

forked and modified from https://github.com/SaturnsVoid/HWID-Based-License-System

## Usage

### Buld and run server

```bash
docker build -t licserver . 
docker run -d --rm -p 9347:9347 -v $(pwd)/db:/app/db --name licserver -e SALT="12345salt" licserver 
```

### Add first license
```bash
 docker exec -it licserver /app/server                                                               
License Server
Github: https://github.com/alexshnup/easy-license-system
Forked from: https://github.com/SaturnsVoid/HWID-Based-License-System
Total Licenses: 0
 
$> add
License Email: test
License Experation (YYYY-MM-DD): 2050-01-01
New License Generated: L83E-ASVL-8MMN for test
$> exit
```

### Run client first time with activate License
```bash
License=L83E-ASVL-8MMN go run client/client.go
...
2023/10/02 22:15:03 Block: block storage (4 disks, 2TB physical storage), Disk Serial: 0ba018e2c3b10023
2023/10/02 22:15:03 HwinfoInit OK
LC....
License file not found.
Try activate license from EnvKey: L83E-ASVL-8MMN
HWID: 5478a74300bb05e53d82f45b4807284c
Connecting to license server...
Registered!
DO NOT DELETE THE FILE! license.dat
 
License OK
```


### Next starts
```bash
go run client/client.go
...
2023/10/02 22:17:04 Block: block storage (4 disks, 2TB physical storage), Disk Serial: 0ba018e2c3b10023
2023/10/02 22:17:04 HwinfoInit OK
LC....
HWID: 5478a74300bb05e53d82f45b4807284c
License OK
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
