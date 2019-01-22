package tests

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	ravendb "github.com/ravendb/ravendb-go-client"
)

func NewTestServiceLocator() (*RavenServerLocator, error) {
	locator, err := NewRavenServerLocator()
	if err != nil {
		return nil, err
	}
	locator.commandArguments = []string{
		"--ServerUrl=http://127.0.0.1:0",
		"--ServerUrl.Tcp=tcp://127.0.0.1:38881",
	}
	return locator, nil
}

const (
	envHTTPSServerURL = "RAVENDB_JAVA_TEST_HTTPS_SERVER_URL"
)

func NewSecuredServiceLocator() (*RavenServerLocator, error) {
	locator, err := NewRavenServerLocator()
	if err != nil {
		return nil, err
	}
	httpsServerURL := os.Getenv(envHTTPSServerURL)
	if httpsServerURL == "" {
		return nil, fmt.Errorf("Unable to find RavenDB https server url. Please make sure %s environment variable is set and is valid (current value = %v)", envHTTPSServerURL, envHTTPSServerURL)
	}

	envCertificatePath := "RAVENDB_JAVA_TEST_CERTIFICATE_PATH"
	certificatePath := os.Getenv(envCertificatePath)
	if certificatePath == "" {
		return nil, fmt.Errorf("Unable to find RavenDB server certificate path. Please make sure %s environment variable is set and is valid (current value = %v)", envCertificatePath, envHTTPSServerURL)
	}

	parsed, err := url.Parse(httpsServerURL)
	if err != nil {
		fmt.Printf("url.Parse('%s') failed with %s\n", httpsServerURL, err)
		return nil, err
	}
	host := parsed.Host
	// host can be name:port, extract "name" part
	host = strings.Split(host, ":")[0]
	tcpServerURL := "tcp://" + host + ":38882"

	locator.commandArguments = []string{
		"--Security.Certificate.Path=" + certificatePath,
		"--Security.Certificate.Password=pwd1234",
		"--ServerUrl=" + httpsServerURL,
		"--ServerUrl.Tcp=" + tcpServerURL,
	}
	return locator, nil
}

func setupRevisions(store *ravendb.DocumentStore, purgeOnDelete bool, minimumRevisionsToKeep int) (*ravendb.ConfigureRevisionsOperationResult, error) {

	revisionsConfiguration := &ravendb.RevisionsConfiguration{}
	defaultCollection := &ravendb.RevisionsCollectionConfiguration{}
	defaultCollection.PurgeOnDelete = purgeOnDelete
	defaultCollection.MinimumRevisionsToKeep = minimumRevisionsToKeep

	revisionsConfiguration.DefaultConfig = defaultCollection
	operation := ravendb.NewConfigureRevisionsOperation(revisionsConfiguration)

	err := store.Maintenance().Send(operation)
	if err != nil {
		return nil, err
	}

	return operation.Command.Result, nil
}
