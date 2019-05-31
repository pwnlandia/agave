package agave

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/d1str0/hpfeeds"
)

const (
	TestHost    = "127.0.0.1"
	TestPort    = 10001
	TestIdent   = "test_id"
	TestAuth    = "test_pass"
	TestChannel = "test_chan"
)

func TestMain(m *testing.M) {
	go startTestBroker()
	time.Sleep(1 * time.Second)

	os.Exit(m.Run())
}

func TestNewHpfeedsWriter(t *testing.T) {
	w, err := NewHpfeedsWriter(TestHost, TestPort, TestIdent, TestAuth, TestChannel)
	if err != nil {
		t.Errorf("Unexpected error building new HpfeedsWriter: %s", err.Error())
	}
	if w == nil {
		t.Fatalf("Nil HpfeedsWriter returned")
	}
}

func startTestBroker() {
	db := NewTestDB()
	b := &hpfeeds.Broker{
		Name: "AgaveTestBroker",
		Port: TestPort,
		DB:   db,
	}

	b.SetDebugLogger(log.Print)
	b.SetInfoLogger(log.Print)
	b.SetErrorLogger(log.Print)

	err := b.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type TestDB struct {
	IDs []hpfeeds.Identity
}

func NewTestDB() *TestDB {
	i := hpfeeds.Identity{
		Ident:       TestIdent,
		Secret:      TestAuth,
		SubChannels: []string{TestChannel},
		PubChannels: []string{TestChannel},
	}
	t := &TestDB{IDs: []hpfeeds.Identity{i}}
	return t
}

func (t *TestDB) Identify(ident string) (*hpfeeds.Identity, error) {
	if ident == TestIdent {
		return &t.IDs[0], nil
	}
	return nil, errors.New("identifier: Unknown identity")
}
