package functions

import "math"

// Helper functions
func max(values ...int) int {
	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func calculateMetrics(actual, predicted []float64) (mae, mse, rmse float64) {
	n := len(actual)
	if n == 0 || n != len(predicted) {
		return 0, 0, 0
	}

	for i := 0; i < n; i++ {
		diff := actual[i] - predicted[i]
		mae += math.Abs(diff)
		mse += diff * diff
	}

	mae /= float64(n)
	mse /= float64(n)
	rmse = math.Sqrt(mse)
	return
}
