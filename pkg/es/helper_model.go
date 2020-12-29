package es

import "encoding/json"

type SearchDocsResponse struct {
	Took         uint64                          `json:"took,omitempty"`
	TimedOut     bool                            `json:"timed_out,omitempty"`
	Shards       *SearchDocsResponseShards       `json:"_shards,omitempty"`
	Hits         *SearchDocsResponseHitsHits     `json:"hits,omitempty"`
	Aggregations *SearchDocsResponseAggregations `json:"aggregations,omitempty"`
}

type SearchDocsResponseShards struct {
	Total      int `json:"total,omitempty"`
	Successful int `json:"successful,omitempty"`
	Skipped    int `json:"skipped,omitempty"`
	Failed     int `json:"failed,omitempty"`
}

type SearchDocsResponseHitsHits struct {
	Total    SearchDocsResponseHitsTotal      `json:"total,omitempty"`
	MaxScore float32                          `json:"max_score,omitempty"`
	Hits     []SearchDocsResponseHitsHitsItem `json:"hits,omitempty"`
}

type SearchDocsResponseHitsTotal struct {
	Value    uint64 `json:"value,omitempty"`
	Relation string `json:"relation,omitempty"`
}

type SearchDocsResponseHitsHitsItem struct {
	Index  string          `json:"_index,omitempty"`
	Type   string          `json:"_type,omitempty"`
	ID     string          `json:"_id,omitempty"`
	Score  float32         `json:"_score,omitempty"`
	Source json.RawMessage `json:"_source,omitempty"`
}

type SearchDocsResponseAggregations struct {
	Result *SearchDocsResponseAggregationsResult `json:"result,omitempty"`
}

type SearchDocsResponseAggregationsResult struct {
	Buckets []SearchDocsResponseAggregationsResultBucket `json:"buckets,omitempty"`
}

type SearchDocsResponseAggregationsResultBucket struct {
	Key      string `json:"key,omitempty"`
	DocCount int    `json:"doc_count,omitempty"`
}
