package obj

type (
	FieldObj struct {
		Name       string
		Type       string
		Embedded   *StructTypeObj
		Visibility bool
	}

	// Obtained after traversing the array of `*ast.Field`
	extractedFieldsData struct {
		usedPackages []UsedPackage
		fieldsSet    []*FieldObj
	}
)
