## Setting Up and Running the Application

Before running the application, make sure you have Docker and Docker Compose installed on your system. 
Additionally, install `golang-migrate` using the following command:

```bash
brew install golang-migrate
```

Next, run Temporal locally by following the instructions at [Temporal Docker Compose](https://github.com/temporalio/docker-compose).

Once Temporal is set up, execute the following commands to prepare and run the application:

```bash
make prepare
```

This command creates the necessary database, runs migrations, and seeds the required data.

```bash
make run
```
This command uses Docker Compose to start both services required for the application. 
After running these commands, your application should be up and running locally.
