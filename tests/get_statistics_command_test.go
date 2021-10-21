package tests

import (
	"testing"

	"github.com/ravendb/ravendb-go-client"
	"github.com/stretchr/testify/assert"
)

func getStatisticsCommandTestCanGetStats(t *testing.T, driver *RavenTestDriver) {
	var err error
	store := driver.getDocumentStoreMust(t)
	defer store.Close()

	executor := store.GetRequestExecutor("")

	sampleData := ravendb.NewCreateSampleDataOperation()
	err = store.Maintenance().Send(sampleData)
	assert.NoError(t, err)

	err = driver.waitForIndexing(store, store.GetDatabase(), 0)
	assert.NoError(t, err)
	command := ravendb.NewGetStatisticsCommand("")
	assert.NoError(t, err)
	err = executor.ExecuteCommand(command, nil)
	assert.NoError(t, err)
	stats := command.Result
	assert.NotNil(t, stats)
	assert.True(t, stats.LastDocEtag > 0)
	assert.Equal(t, stats.CountOfIndexes, 7)
	assert.Equal(t, stats.CountOfDocuments, int64(1059))
	assert.True(t, stats.CountOfRevisionDocuments > 0)
	assert.Equal(t, stats.CountOfDocumentsConflicts, int64(0))
	assert.Equal(t, stats.CountOfConflicts, int64(0))
	assert.Equal(t, stats.CountOfUniqueAttachments, int64(17))
	assert.NotEqual(t, stats.DatabaseChangeVector, "")
	assert.NotEqual(t, stats.DatabaseID, "")
	assert.NotNil(t, stats.Pager)
	assert.NotNil(t, stats.GetLastIndexingTime())
	assert.NotNil(t, stats.Indexes)
	assert.NotEqual(t, stats.SizeOnDisk.HumaneSize, "")
	assert.NotEqual(t, stats.SizeOnDisk.SizeInBytes, 0)

	indexes := stats.Indexes
	for _, indexInformation := range indexes {
		assert.NotEqual(t, indexInformation.Name, "")
		assert.False(t, indexInformation.IsStale)
		assert.NotNil(t, indexInformation.State)
		assert.NotEqual(t, indexInformation.LockMode, "")
		assert.NotEqual(t, indexInformation.Priority, "")
		assert.NotEqual(t, indexInformation.Type, "")
		assert.NotNil(t, indexInformation.GetLastIndexingTime())
	}
}

func TestGetStatisticsCommand(t *testing.T) {
	driver := createTestDriver(t)
	destroy := func() { destroyDriver(t, driver) }
	defer recoverTest(t, destroy)

	getStatisticsCommandTestCanGetStats(t, driver)
}
