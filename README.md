# Foosball tournament [![Go Report Card](https://goreportcard.com/badge/github.com/vlarrat-theodo/lbc-foosball)](https://goreportcard.com/report/github.com/vlarrat-theodo/lbc-foosball)

Foosball is a game between two users, particularly liked in startups.

Each user has a team of 11 players (named p1 - p2 - p3 ... from goalkeeper to forwards).
```
                    +----------+
                    |          |
    +------------------------------------------+
    |                                          |
    |                    p1                    | # goalkeeper
    |                                          |
    |                                          |
    |            p2              p3            | # defenders
    |                                          |
    |                                          |
    |      p4     p5     p6     p7     p8      | # midfielders
    |                                          |
    |                                          |
    |        p9         p10         p11        | # forwards
    |                                          |
    +------------------------------------------+
```

Each of the players can score goals, the scoring depends on:
- which player scored
- whether the goal was a "gamelle" or not
The exact rules are specified in the "Rules" section.
This API stores all tournament's results:
- a user can have multiple opponents
- for each couple of users, there will be at most one current score

### Goal
Provide an API that exposes a POST route to store a new tournament goal
and returns the current score between these two users.
As a simplification, it is recommended not to use an external database.
- POST /goal:
- the user who scored (represented by his ID), ex. user1
- his opponent (represented by his ID), ex. user2
- the field position of a scorer (ex. "p3")
- whether the goal is a "gamelle" or not (true for a gamelle)


### Bonus
Expose a GET route to know a user's set balance in total, against all users.

The sets must be finished to be counted.
```
GET /balance?user_id=<user_id>
Returns: {"won": 5, "lost": 3}
```


### Additional information
[Foosball rules](./docs/foosball_rules.md)

[Foosball examples](./docs/foosball_examples.md)

---

## Project information
This project has been developed using [GO language](https://golang.org/) and is using serverless [Lambda](https://aws.amazon.com/en/lambda/features/) architecture.

It has been deployed online on AWS resources (RDS database, Lambda functions and API Gateway) and can be tested on this architecture.
It can also be installed and tested locally (in this case, AWS resources are emulated).

---

## Test project online
Online architecture is available at [https://api.lbc-foosball.theo.do](https://api.lbc-foosball.theo.do).

To test goal submission route, use following cURL command (adapt body content to your expectations):
```shell script
curl -X POST \
  https://api.lbc-foosball.theo.do/goal \
  -d '{
	"scorer": "user1",
	"opponent": "user2",
	"player": "p1",
	"gamelle": false
}'
```

To test user set balance route, use following cURL command (adapt query parameter to your expectations):
```shell script
curl -X GET \
  'https://api.lbc-foosball.theo.do/balance?user_id=user1'
```

---

## Install project locally
#### Prerequisites
To make this project work on your machine, you need to install following prerequisites:
- [GO](https://golang.org/dl/) version 1.12.6 or above
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html) to manage AWS resources
- [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html) to emulate AWS resources locally

Once finished, just go to main folder and launch following command:
```shell script
make install
```
It will install the complete project and its prerequisites.

For more available commands, just launch:
```shell script
make help
```

---

## Run and test project locally
Once installed, you can launch local API by launching following command:
```shell script
make start
```
Local API is available at [http://localhost:3000](http://localhost:3000).
                      
To test goal submission route, use following cURL command (adapt body content to your expectations):
```shell script
curl -X POST \
 http://localhost:3000/goal \
 -d '{
"scorer": "user1",
"opponent": "user2",
"player": "p1",
"gamelle": false
}'
```

To test user set balance route, use following cURL command (adapt query parameter to your expectations):
```shell script
curl -X GET \
 'http://localhost:3000/balance?user_id=user1'
```

To launch GO tests with coverage, just launch following command:
```shell script
make test
```
