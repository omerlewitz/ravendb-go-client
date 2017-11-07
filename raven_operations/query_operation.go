package raven_operations

import (
	"github.com/ravendb-go-client/http/commands"
	"errors"
	"github.com/ravendb-go-client/store"
	"github.com/ravendb-go-client/http/server_nodes"
	"fmt"
	"net/http"
	"github.com/ravendb-go-client/tools"
)
// TODO: move into ./data/query.go
//@param str query: Actual query that will be performed.
//@param dict query_parameters: Parameters to the query
//@param int start:  Number of records that should be skipped.
//@param int page_size:  Maximum number of records that will be retrieved.
//@param bool wait_for_non_stale_results: True to wait in the server side to non stale result
//@param None or float cutoff_etag: Gets or sets the cutoff etag.

type IndexQuery struct {
	PageSize uint
	PageSizeSet bool
}
func (ref * IndexQuery) get_query_hash() string{
	return "hash"
}
func (obj IndexQuery) to_json() interface{} {
	return ""
}
//@param IndexQuery indexQuery: A query definition containing all information required to query a specified index.
//@param bool metadataOnly: True if returned documents should include only metadata without a document body.
//@param bool indexEntriesOnly: True if query results should contain only index entries.
//@return:json
type QueryOperation struct {
	commands.Command
	session  *store.DocumentSession
	indexName string
	indexQuery IndexQuery
	disableEntitiesTracking, metadataOnly, indexEntriesOnly bool
}
func NewQueryOperation(session *store.DocumentSession, indexName string, indexQuery IndexQuery, disableEntitiesTracking, metadataOnly, indexEntriesOnly bool) (*QueryOperation, error) {

	if session.GetConvention().RaiseIfQueryPageSizeIsNotSet && !indexQuery.PageSizeSet {
		return nil, errors.New(`Attempt to query without explicitly specifying a page size.
		You can use .take() methods to set maximum number of results.
		By default the page size is set to sys.maxsize and can cause
		severe performance degradation.`)
	}
	ref := &QueryOperation{}
	ref.session = session
	ref.indexName = indexName
	ref.indexQuery = indexQuery
	ref.disableEntitiesTracking = disableEntitiesTracking
	ref.metadataOnly = metadataOnly
	ref.indexEntriesOnly = indexEntriesOnly

	return ref, nil

}
func (ref *QueryOperation) GetIndexQuery() IndexQuery{
	return ref.indexQuery
}
func (ref *QueryOperation) createRequest(serverNode server_nodes.IServerNode) {
	ref.session.IncrementRequestsCount()
	// will implement logging later
	//logging.info("Executing query '{0}' on index '{1}'".format(ref._index_query.query, ref._index_name))

	ref.Url = fmt.Sprintf("%s/databases/%s/queries?query-hash=%s", serverNode.GetUrl(), serverNode.GetDatabase(), ref.indexQuery.get_query_hash())
	if ref.metadataOnly {
		ref.Url += "&metadata-only=true"
	}
	if ref.indexEntriesOnly {
		ref.Url += "&debug=entries"
	}

	ref.Data = ref.indexQuery.to_json()
}
func (ref *QueryOperation) GetResponseRaw(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, errors.New("response is nil")
	}

	return tools.ResponseToJSON(resp)
}