package analyze

import "github.com/g10z3r/archx/internal/analyze/types"

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

func CalculateLCOM(s *types.StructType) float32 {
	if s == nil || len(s.Methods) == 0 || len(s.Fields) == 0 {
		return defaultLCOMValue
	}

	P := 0
	Q := 0

	for fieldName := range s.Fields {
		for _, fieldSet := range s.Methods {
			// Check if the current method uses the current field
			if _, exists := fieldSet[fieldName]; exists {
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

func CalculateLCOM96B(s *types.StructType) float32 {
	if s == nil || len(s.Methods) == 0 || len(s.Fields) == 0 {
		return defaultLCOMValue
	}

	// Total number of attributes in the class
	M := float32(len(s.Fields))
	// Number of methods in the class
	n := float32(len(s.Methods))

	var sum float32 = 0.0
	for _, fields := range s.Methods {
		// The number of attributes that method i works with
		m_i := len(fields)
		sum += 1.0 - (float32(m_i) / float32(M))
	}

	return sum / float32(n)
}