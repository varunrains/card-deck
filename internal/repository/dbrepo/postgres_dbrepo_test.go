package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/varunrains/carddeck/internal/repository"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "cardDB_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo repository.DatabaseRepo
var deckId string

func TestMain(m *testing.M) {
	//connect to docker; fail if docker not running
	p, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	pool = p

	//set up our docker options, specifying the image and so forth
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	//get a resource (docker image)
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resources: %s", err)
	}

	//start the image and wait until it's ready
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database:%s", err)
	}

	//populate the database with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepo = &PostgresDBRepo{DB: testDB}
	//run tests
	code := m.Run()

	//clean up
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/create_tables.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

func TestCreateDeck(t *testing.T) {
	deck, err := testRepo.CreateDeck(false, []string{})
	deckId = deck.ID
	if err != nil {
		t.Errorf("create deck returned an error: %s", err)
	}

	if deck.ID != deckId {
		t.Errorf("create deck didn't create sucessfully; expected %s but got %s", deckId, deck.ID)
	}
}

func TestCreateDeckWithShuffledAndSpecifiedCards(t *testing.T) {
	cards := []string{"AS", "KD", "AC", "2C", "KH"}
	var count = 0
	deck, err := testRepo.CreateDeck(true, cards)

	if err != nil {
		t.Errorf("create deck returned an error: %s", err)
	}

	if len(cards) != len(deck.Cards) {
		t.Errorf("remaining cards should have been %d but got %d", len(cards), len(deck.Cards))
	}

	for i, c := range deck.Cards {
		c.SetCode()
		if c.Code == cards[i] {
			count++
		}
	}

	if count == 5 {
		t.Errorf("cards are not shuffled")
	}
}

func TestCreateDeckWithFullCards(t *testing.T) {
	deck, err := testRepo.CreateDeck(false, []string{})
	if err != nil {
		t.Errorf("create deck returned an error: %s", err)
	}

	if len(deck.Cards) != 52 {
		t.Errorf("remaining cards should have been %d but got %d", 52, len(deck.Cards))
	}
}

func TestDrawDeck(t *testing.T) {
	cards, err := testRepo.DrawDeck(deckId, 1)
	if err != nil {
		t.Errorf("draw deck returned an error: %s", err)
	}

	if len(*cards) != 1 {
		t.Errorf("remaining cards after drawing is %d but expected %d", len(*cards), 1)
	}
}

func TestDrawDeckWithMoreNumberThanExisting(t *testing.T) {
	var errorvalue = "remaining cards in the deck is less than the count"
	deck, err := testRepo.CreateDeck(false, []string{})

	if err != nil {
		t.Errorf("create deck returned an error: %s", err)
	}

	_, err = testRepo.DrawDeck(deck.ID, 53)

	if err.Error() != errorvalue {
		t.Errorf("error while drawing cards; expected %s but got %s", errorvalue, err.Error())
	}
}

func TestOpenDeck(t *testing.T) {
	deck, err := testRepo.OpenDeck(deckId)
	if err != nil {
		t.Errorf("open deck returned an error: %s", err)
	}

	if deck.ID != deckId || len(deck.Cards) != 51 {
		t.Errorf("number of cards after opening from a deck is wrong; expected %d but got %d", 1, len(deck.Cards))
	}
}

func TestOpenDeckWithInvalidId(t *testing.T) {
	_, err := testRepo.OpenDeck(uuid.NewString())
	var expectedError = "not a valid deck id/ deck is empty"
	if err.Error() != expectedError {
		t.Errorf("expected error %s: but got %s", expectedError, err)
	}
}

func TestDrawDeckWithInvalidId(t *testing.T) {
	_, err := testRepo.DrawDeck(uuid.NewString(), 1)
	var expectedError = "not a valid deck id"
	if err.Error() != expectedError {
		t.Errorf("expected error %s: but got %s", expectedError, err)
	}
}
