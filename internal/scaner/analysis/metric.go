package scaner

import "github.com/g10z3r/archx/internal/scaner/entity"

const (
	defaultLCOMValue         = 0.0
	defaultAbstractnessValue = 0.0
)

func CalculateAbstractness(abstractEntities int, specificEntities int) float32 {
	if abstractEntities == 0 && specificEntities == 0 {
		return defaultAbstractnessValue
	}

	return float32(abstractEntities) / float32(abstractEntities+specificEntities)
}

func CalculateLCOM(s *entity.StructInfo) float32 {
	if s == nil || len(s.Methods) == 0 || len(s.Fields) == 0 {
		return defaultLCOMValue
	}

	P := 0
	Q := 0

	for fieldName := range s.FieldsIndex {
		for _, methodInfo := range s.Methods {
			// Check if the current method uses the current field
			if _, exists := methodInfo.UsedFields[fieldName]; exists {
				Q++
			} else {
				P++
			}
		}
	}

	if P > Q {
		return float32(P - Q)
	}

	return defaultLCOMValue
}

func CalculateLCOM96B(s *entity.StructInfo) float32 {
	if s == nil || len(s.Methods) == 0 || len(s.Fields) == 0 {
		return defaultLCOMValue
	}

	// Total number of attributes in the class
	M := float32(len(s.Fields))
	// Number of methods in the class
	n := float32(len(s.Methods))

	var sum float32 = 0.0
	for _, method := range s.Methods {
		// The number of attributes that method i works with
		m_i := len(method.UsedFields)
		sum += 1.0 - (float32(m_i) / M)
	}

	return sum / n
}
