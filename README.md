# Extended-HWID-Based-License-System

An evolved GoLANG based License System, extending the original basic HWID (hardware ID) license system into a comprehensive and API-driven management tool.

## Acknowledgment and Modifications

This project has been forked and significantly modified from the original [HWID-Based-License-System](https://github.com/SaturnsVoid/HWID-Based-License-System) developed by [SaturnsVoid](https://github.com/SaturnsVoid). I deeply appreciate the initial work and foundation laid by the original developer.

### Major Changes from the Original Version

- **Extended API Functionality**: 
    - The system now comes with an extensive HTTP API, allowing users to interact with the license system remotely and programmatically. This allows operations such as adding, listing, resetting, and removing keys through HTTP requests.
  
- **Improved Security**: 
    - Transitioned from MD5 to SHA-256 for hashing, providing a more secure and collision-resistant hashing algorithm to safeguard the integrity and security of data.
    - The key point was that the response from the server comes in the form of a hash and is checked on the client side
    - The security of the system has been enhanced by introducing token-based authorization for API access, ensuring that only authorized users can interact with critical API endpoints.

- **Containerization**: 
    - Docker support has been implemented, allowing both the client and server to be containerized, which enhances the portability and deployment flexibility of the system.

While I have maintained the essence and usability of the original system, my version introduces several pivotal changes and enhancements aimed at providing additional functionality, improving user experience, and ensuring a higher level of security and stability.


## Overview

While the original concept sparked the creation of this advanced version, the codebase has been substantially reconstructed to cater to more varied and complex use-cases, pivoting from a simplistic HWID verification system to an API-oriented, robust license management system.


## Usage

Utilize Docker to encapsulate and deploy the license server.

### Build and Run Server
```bash
docker network create lic # create network for server and client to communicate inside Docker for testing
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

### Run Client
```bash
License=86UU-N4SB-OQYH go run client/client.go 
```

### Run Client in Docker

first launch of the client with license activation
```bash
docker build -t client -f Dockerfile-client . 
docker run --network lic -t --name client -e License="16FB-L6AX-2ZPZ" client
```

next client launch with license check
```bash
docker start client
docker logs client | tail -n 2
```

## Key Enhancements

- **API-Oriented Interactions**: Now manage and interact with license data programmatically via HTTP API calls.
  - **Add Key**: API endpoint to add a key to the license database.
    ```bash
        curl -X POST http://127.0.0.1:9347/add \
        -H "Authorization: mytoken" \
        -H "Content-Type: application/json" \
        -d '{"email": "test2", "expiration": "2023-12-31"}'

        {"email":"test2","exp_date":"2023-12-31","license":"5JDG-DVFC-5Z3H","message":"New license generated"}
    ```
  - **List Keys**: Accessible API endpoint to list all keys.
    ```bash
        curl -s -X GET http://127.0.0.1:9347/list \
        -H "Authorization: mytoken" | jq

        [
        "16FB-L6AX-2ZPZ:2023-12-31:test1:8e4db87551f6decab29760f2c2c0b8a74a0b746f08805f035cdb54c0923b4db5",
        "5JDG-DVFC-5Z3H:2023-12-31:test2:NOTSET"
        ]
    ```
  - **Reset Key**: API feature to reset keys, altering the status within the database.
    ```bash
        curl -X POST http://127.0.0.1:9347/reset-key \
        -H "Authorization: mytoken" \
        -H "Content-Type: application/json" \
        -d '{"key": "16FB-L6AX-2ZPZ"}'
            
        {"status":"success"}
    ```
  - **Remove Key**: Remove keys from the database via API calls.
    ```bash
        curl -X DELETE http://127.0.0.1:9347/remove \
        -H "Authorization: mytoken" \
        -H "Content-Type: application/json" \
        -d '{"email": "test1"}'

        {"status":"success"}
    ```
- **Token-Based Authentication**: Introduce a secure layer for API interactions (**except for 'Check method'**).

## Workflow

### 1. Key Generation:
- Generate keys using the API or manually, associating them with Email and Expiration Date.
- Keys conform to a 4x4x4 character format, randomly generated from characters 0-9 and A-Z.

### 2. Client Interaction:
- Initial run searches for `license.dat`, prompting client registration if not found.
- Clients input their key, subsequently generating a unique HWID for their system/user.
- The software verifies key validity against the license server, ensuring it's not expired or linked to another HWID.
- Successful validation adds the HWID to the database.

### 3. Key Management:
- Execute key management tasks through API calls (add, list, reset, remove).
- Ensure secure API interactions using token-based authentication (excluding 'List Keys').

### 4. License Verification:
- The software can initiate a license check against the server as required. (TODO: implement this feature)
- License verification can be spontaneous or operate on a timed loop.

### 5. Server Flexibility:
- The license server can be configured to run on any available port, adapting to diverse networking scenarios.

## Web Innterface
- (TODO: implement this feature)

## Security & Usage Note:

- **Security Consideration**: Use API endpoints judiciously, ensuring secure communication (preferably HTTPS) to prevent inadvertent data exposure.
- **Database Management**: The system utilizes a text file for the database but offers structured manipulation via API calls, minimizing manual interactions.

## Disclaimer:

This system, although providing an advanced level of license management, is not impenetrable. It is intended to serve as a sophisticated deterrent, offering secure and manageable functionalities. Use it with a comprehensive understanding of its capabilities and boundaries.

