package es

import "encoding/json"

type InfoResp struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
	Version     InfoRespVersion
}

type InfoRespVersion struct {
	Number                           string `json:"number"`
	BuildFlavor                      string `json:"build_flavor"`
	BuildType                        string `json:"build_type"`
	BuildHash                        string `json:"build_hash"`
	BuildDate                        string `json:"build_date"`
	BuildSnapshot                    bool   `json:"build_snapshot"`
	LuceneVersion                    string `json:"lucene_version"`
	MinimumWireCompatibilityVersion  string `json:"minimum_wire_compatibility_version"`
	MinimumIndexCompatibilityVersion string `json:"minimum_index_compatibility_version"`
}

type ListDocsResp struct {
	Took         uint64                    `json:"took,omitempty"`
	TimedOut     bool                      `json:"timed_out,omitempty"`
	Shards       *ListDocsRespShards       `json:"_shards,omitempty"`
	Hits         *ListDocsRespHitsHits     `json:"hits,omitempty"`
	Aggregations *ListDocsRespAggregations `json:"aggregations,omitempty"`
}

type ListDocsRespShards struct {
	Total      int `json:"total,omitempty"`
	Successful int `json:"successful,omitempty"`
	Skipped    int `json:"skipped,omitempty"`
	Failed     int `json:"failed,omitempty"`
}

type ListDocsRespHitsHits struct {
	Total    ListDocsRespHitsTotal      `json:"total,omitempty"`
	MaxScore float32                    `json:"max_score,omitempty"`
	Hits     []ListDocsRespHitsHitsItem `json:"hits,omitempty"`
}

type ListDocsRespHitsTotal struct {
	Value    uint64 `json:"value,omitempty"`
	Relation string `json:"relation,omitempty"`
}

type ListDocsRespHitsHitsItem struct {
	Index  string          `json:"_index,omitempty"`
	Type   string          `json:"_type,omitempty"`
	ID     string          `json:"_id,omitempty"`
	Score  float32         `json:"_score,omitempty"`
	Source json.RawMessage `json:"_source,omitempty"`
}

type ListDocsRespAggregations struct {
	Result *ListDocsRespAggregationsResult `json:"result,omitempty"`
}

type ListDocsRespAggregationsResult struct {
	Buckets []ListDocsRespAggregationsResultBucket `json:"buckets,omitempty"`
}

type ListDocsRespAggregationsResultBucket struct {
	Key      string `json:"key,omitempty"`
	DocCount int    `json:"doc_count,omitempty"`
}

type ListIndicesResp struct {
	Items []ListIndicesItemResp
}

type ListIndicesItemResp struct {
	ID           string `json:"uuid"`
	Name         string `json:"index"`
	Health       string `json:"health"`
	Status       string `json:"status"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

type CreateIndexResp struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	Index              string `json:"index"`
}

type DeleteIndexResp struct {
	Acknowledged bool `json:"acknowledged"`
}

type ReindexResp struct {
	Took                 uint64             `json:"took"`
	TimedOut             bool               `json:"timed_out"`
	Total                uint64             `json:"total"`
	Updated              uint64             `json:"updated"`
	Created              uint64             `json:"created"`
	Deleted              uint64             `json:"deleted"`
	Batches              uint64             `json:"batches"`
	VersionConflicts     uint64             `json:"version_conflicts"`
	Noops                uint64             `json:"noops"`
	Retries              ReindexRespRetries `json:"retries"`
	ThrottledMillis      uint64             `json:"throttled_millis"`
	RequestsPerSecond    float32            `json:"requests_per_second"`
	ThrottledUntilMillis uint64             `json:"throttled_until_millis"`
	Failures             []interface{}      `json:"failures"`
}

type ReindexRespRetries struct {
	Bulk   uint64 `json:"bulk"`
	Search uint64 `json:"search"`
}

type AliasOrUnaliasIndexResp struct {
	Acknowledged bool `json:"acknowledged"`
}

type ListAliasesResp struct {
	Items []ListAliasesRespItem
}

type ListAliasesRespItem struct {
	Alias         string `json:"alias"`
	Index         string `json:"index"`
	Filter        string `json:"filter"`
	RoutingIndex  string `json:"routing.index"`
	RoutingSearch string `json:"routing.search"`
}
