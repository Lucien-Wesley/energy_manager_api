package functions

// difference performs differencing on a time series
func difference(series []float64, order int) []float64 {
	if order <= 0 {
		return series
	}
	result := make([]float64, len(series)-order)
	for i := order; i < len(series); i++ {
		result[i-order] = series[i] - series[i-order]
	}
	return result
}

// inverseDifference reverses the differencing operation
func inverseDifference(diff []float64, original []float64, order int) []float64 {
	result := make([]float64, len(diff))
	copy(result, diff)

	for i := 0; i < len(result); i++ {
		if i-order >= 0 {
			result[i] += original[i]
		} else if i < len(original) {
			result[i] += original[i]
		}
	}
	return result
}