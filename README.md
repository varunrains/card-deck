# card-deck

#### Installation
To install this project, first clone the repository:

git clone `https://github.com/varunrains/carddeck.git`

#### Build

`make build` - Builds all projects

#### Run locally

`make run` - Runs the API [Please use docker desktop](https://www.docker.com/products/docker-desktop/)

#### To RUN Tests

`make test_cmd_api` To run the tests under ./cmd/api (Routes and Handlers)
<br/><br/>
`make test_integration` To run the DB integration tests.
<br/><br/>
`make test_coverage` To check the test coverage

#### How to test the API's ?

1) For Creating the Deck use : `http://localhost:6555/createDeck?shuffled=true&cards=AH,AD,AS,AC` this is a **POST** request.<br/>
**Query Parameters:** <br/> `shuffled` To randomize the cards in the deck while creation.  <br/> `cards` To specify the cards during creation. Eg: `AH,AD` stands for **ACE of HEARTS** and **ACE of DIAMONDS** <br /><br/>
2) For Opening the Deck use : `http://localhost:6555/openDeck/81729356-c7ec-4f87-b0ae-bc08e4800b2a` this is a **GET** request. <br/><br/>
3) For Drawing the Deck use : `http://localhost:6555/drawDeck/81729356-c7ec-4f87-b0ae-bc08e4800b2a?count=4` this is a **GET** request. <br/>
**Query Parameters:** <br/> `count` Number of cards to draw from the deck.

#### Note:
1) With docker desktop the application runs in port `6555`.
2) Database used is `PostgreSQL` and it will be installed through `Docker Desktop`.
3) All the steps for the installation of the `micro-service` (carddeck) and `PostgreSQL` is present in the `docker-compose.yml` file.
4)`Makefile` is used for building and running the API using Docker desktop and also to run the tests.