package utils

func GetFloat64Value(value string) float64 {
	switch value {
	case "A":
		return 100
	case "B":
		return 85
	case "C":
		return 70
	case "D":
		return 60
	case "E":
		return 50
	}

	return 50
}
