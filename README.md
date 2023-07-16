# Federation Hub Setup

This repository contains the code for setting up a federated marketplace using Go and PostgreSQL.

## Prerequisites

- Go
- PostgreSQL

## Getting Started

1. Clone the repository:

   `git clone https://github.com/ADSP-Project/Federation-Hub.git`

   `cd Federation-Hub`

2. Set up the PostgreSQL database:
   - Install PostgreSQL and connect:
     
     `sudo apt install postgresql postgresql-contrib`

     `psql`

   - Create a new database named 'shops':

        `CREATE DATABASE shops;`

   - Create a new user with appropriate privileges:

        `CREATE USER your_username WITH PASSWORD 'your_password';`

        `GRANT ALL PRIVILEGES ON DATABASE shops TO your_username;`

   - Connect to Postgres with your new user:
     
     `psql -d shops -U your_username`

   - Create a 'shops' table in the 'shops' database with columns 'id', 'name', and 'webhookURL':

        ``CREATE TABLE shops (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255),
        description VARCHAR(255),
        webhookURL VARCHAR(255),
        publicKey VARCHAR(1024)
        );``

3. Configure environment variables:
   - Rename the `.env.example` file to `.env`.
   - Update the database credentials in the `.env` file to match your PostgreSQL setup.

4. Install dependencies:

   - Make sure to initialize go.mod to manage dependencies:
   
      `go mod init github.com/ADSP-Project/Federation-Hub`

   - Fetch and arrange them into newly generated go.mod:
   
      `go mod tidy`

   - Finally, install the required dependencies with the following command:
   
      `go mod download`

5. Run the main server:
   - Start the Go server using the following command:

        `go run main.go`

   - The server will start running on port 8000 by default.

6. Run the authentification server:
   - Start the authentification server using the following command:

      `go run authServer.go`

   - The server will start running on port 8081 by default.

7. Simulate a shop joining the federation:
   - To simulate a shop joining the federation, open a new terminal and run the following command:

     `go run shop.go [port] [name]`

     Replace `[port]` with the desired port number and `[name]` with the name of the shop.

     This will start a shop server that automatically joins the federation by sending a POST request to the federation server.

8. Access the federated marketplace:
   - You can now access the federated marketplace at `http://localhost:8000`.

9. Additional Notes:
   - You can run multiple instances of the shop server by providing different port numbers and shop names.


## Run with skaffold

   minikube start --profile hub

   skaffold run