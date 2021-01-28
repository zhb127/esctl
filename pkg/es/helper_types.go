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

type CatIndicesResp struct {
	Items []CatIndicesItemResp
}

type CatIndicesItemResp struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	Uuid         string `json:"uuid"`
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
