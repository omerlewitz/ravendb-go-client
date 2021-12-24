package tests

import (
	"github.com/ravendb/ravendb-go-client"
	"github.com/ravendb/ravendb-go-client/serverwide/operations"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func putSecretKey(t *testing.T, driver *RavenTestDriver) {
	var err error
	key := "ET4dgQBqN574qLooWS33dhTM9SiFtJghNiSXP6XnArw="
	if os.Getenv("RAVEN_License") == "" {
		t.Skip("This test requires RavenDB license.")
	}

	store := driver.getSecuredDocumentStoreMust(t)
	assert.NotNil(t, store)
	defer store.Close()

	driver2 := createTestDriver(t)
	store2, err := driver2.createMainStore()
	defer store2.Close()
	assert.NoError(t, err)

	driver3 := createTestDriver(t)
	store3, err := driver3.createMainStore()
	defer store3.Close()
	assert.NoError(t, err)

	destroy := func() { destroyDriver(t, driver) }
	defer recoverTest(t, destroy)

	operationAddNodeToCluster := operations.OperationAddClusterNode{
		Url:     store2.GetUrls()[0],
		Tag:     "B",
		Watcher: false,
	}
	err = store.Maintenance().Server().Send(&operationAddNodeToCluster)
	assert.NoError(t, err)

	operationAddNodeToCluster = operations.OperationAddClusterNode{
		Url:     store3.GetUrls()[0],
		Tag:     "C",
		Watcher: false,
	}
	err = store.Maintenance().Server().Send(&operationAddNodeToCluster)
	assert.NoError(t, err)

	for _, url := range store.GetUrls() {
		name := "test_db"
		command := ravendb.NewPutSecretKeyCommand(name, key, false)
		requestExecutor := ravendb.ClusterRequestExecutorCreateForSingleNode(url, store.Certificate, store.TrustStore, store.GetConventions())
		err := requestExecutor.ExecuteCommand(command, nil)
		if err != nil {
			panic(err)
		}
	}

}

func TestPutSecretKey(t *testing.T) {
	driver := createTestDriver(t)
	destroy := func() {
		destroyDriver(t, driver)
	}

	defer recoverTest(t, destroy)

	putSecretKey(t, driver)
}
