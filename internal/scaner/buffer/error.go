package buffer

import "errors"

var (
	errIncorrectStructBufferType = errors.New("buffer is not of type *StructBuffer")
	errIncorrectImportBufferType = errors.New("buffer is not of type *ImportBuffer")
)
