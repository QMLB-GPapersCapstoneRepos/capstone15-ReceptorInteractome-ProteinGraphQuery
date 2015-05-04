package models

const EntryBucketName string = "graph"
const DatabaseName string = "db/HumanPredictions.db"

type GraphEntry struct {
	Target string
	Score  float64
}

func (e GraphEntry) MakesCutoff(cutoff float64) bool {
	return e.Score >= cutoff
}
