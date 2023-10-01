package dao

type FieldDAO struct {
	Type     string     `bson:"type"`
	Embedded *StructDAO `bson:"embedded"`
	IsPublic bool       `bson:"isPublic"`
}

type MethodDAO struct {
	UsedFields map[string]int
	IsPublic   bool
}

type DependencyDAO struct {
	ImportIndex int `bson:"importIndex"`
	Usage       int `bson:"usage"`
}

type StructDAO struct {
	Fields      []*FieldDAO    `bson:"fields"`
	FieldsIndex map[string]int `bson:"fieldsIndex"`

	Methods      []*MethodDAO   `bson:"methods"`
	MethodsIndex map[string]int `bson:"methodsIndex"`

	Dependencies      []*DependencyDAO `bson:"dependencies"`
	DependenciesIndex map[string]int   `bson:"dependenciesIndex"`

	Incomplete bool `bson:"incomplete"`
	IsEmbedded bool `bson:"isEmbedded"`
}
