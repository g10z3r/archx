package model

type SnapshotDAO struct {
	Timestamp int64  `bson:"timestamp"`
	BasePath  string `bson:"basePath"`
}

// func MapSnapshotEntity(e *entity.SnapshotEntity) *SnapshotDAO {
// 	return &SnapshotDAO{
// 		Timestamp: e.Timestamp,
// 		BasePath:  e.BasePath,
// 	}
// }
