package es

import "encoding/json"

type SearchDocsResp struct {
	Took         uint64                      `json:"took,omitempty"`
	TimedOut     bool                        `json:"timed_out,omitempty"`
	Shards       *SearchDocsRespShards       `json:"_shards,omitempty"`
	Hits         *SearchDocsRespHitsHits     `json:"hits,omitempty"`
	Aggregations *SearchDocsRespAggregations `json:"aggregations,omitempty"`
}

type SearchDocsRespShards struct {
	Total      int `json:"total,omitempty"`
	Successful int `json:"successful,omitempty"`
	Skipped    int `json:"skipped,omitempty"`
	Failed     int `json:"failed,omitempty"`
}

type SearchDocsRespHitsHits struct {
	Total    SearchDocsRespHitsTotal      `json:"total,omitempty"`
	MaxScore float32                      `json:"max_score,omitempty"`
	Hits     []SearchDocsRespHitsHitsItem `json:"hits,omitempty"`
}

type SearchDocsRespHitsTotal struct {
	Value    uint64 `json:"value,omitempty"`
	Relation string `json:"relation,omitempty"`
}

type SearchDocsRespHitsHitsItem struct {
	Index  string          `json:"_index,omitempty"`
	Type   string          `json:"_type,omitempty"`
	ID     string          `json:"_id,omitempty"`
	Score  float32         `json:"_score,omitempty"`
	Source json.RawMessage `json:"_source,omitempty"`
}

type SearchDocsRespAggregations struct {
	Result *SearchDocsRespAggregationsResult `json:"result,omitempty"`
}

type SearchDocsRespAggregationsResult struct {
	Buckets []SearchDocsRespAggregationsResultBucket `json:"buckets,omitempty"`
}

type SearchDocsRespAggregationsResultBucket struct {
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
	Items []ListAliasesItemResp
}

type ListAliasesItemResp struct {
	Alias         string `json:"alias"`
	Index         string `json:"index"`
	Filter        string `json:"filter"`
	RoutingIndex  string `json:"routing.index"`
	RoutingSearch string `json:"routing.search"`
}
