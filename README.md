# payment-system-kafka
_____
## Payment System with Kafka

This project is a transaction system based on Apache Kafka, Redis and PostgreSQL, written in the Go programming language.
The project is divided into three main components: the client, producer, and consumer.

## Overview
The Payment System with Kafka is designed to handle transactions efficiently. It follows a simple flow:

1. `Client`: Sends transaction data to the producer.
2. `Producer`: Publishes data to Kafka topics for further use by the consumer and stores transaction data in Redis.
3. `Consumer`: The consumer subscribes to specific topics in Kafka and processes incoming transaction messages.

## Components
### Client
The client performs the initial stage of data processing in our system.

Tasks:

* Receives or retrieves transaction data to be processed by the system.

* Sends this data to the Kafka producer for further publication.

* Processes user's transactions record and performs authentication based on JWT tokens and Sessions.

### Producer
The producer is a key component between the client and the Kafka system.
It accepts data from the client and publishes it to the designated Kafka topic.
Additionally, the producer stores the transaction data in Redis for reference.

Tasks:

* Receives transaction data from the client.

* Transforms and encrypts transaction data.

* Publishes data to Kafka topics for further use by the consumer.
### Consumer
The consumer subscribes to specific topics in Kafka and processes incoming transaction messages.

Tasks:

* Subscribes to specific Kafka topics where the producer publishes data.

* Receives and processes incoming transaction messages from Kafka.

* Ensures that all messages are processed reliably and in accordance with the system's requirements.

## Client's database
The database schema consists of three relationships: credentials, users, and transactions.

1. `Credentials`:
   Stores user's credentials. Includes fields for a unique identifier (id), unique email (email), password hash (password_hash), and registration date (register_date).
Enforces uniqueness of emails and utilizes a unique constraint on the email field.

2. `Users`:
   Contains user's profile information. Includes fields for a unique identifier (id), first name (first_name), last name (last_name), and a reference to the associated credentials (credential_id). Enforces uniqueness of credential_id and establishes a foreign key constraint (fk_credentials_users) linking to the credentials table.

3. `Transactions`:
   Records transactions made by users. Contains fields for a unique identifier (id), associated user (user_id), transaction identifier (transaction_id), hashed card number (card_number), and transaction amount (amount). Establishes a foreign key constraint (fk_user_transactions) linking to the users table based on the user_id.

![main](https://i.imgur.com/G5uxZ0S.png)
