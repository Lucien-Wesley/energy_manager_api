package functions

import (
	"api/ai_models"
	"math"
	"math/rand"
	"time"
)

// optimizeARIMA optimizes ARIMA coefficients using gradient descent
func optimizeARIMA(series []float64, params aimodels.ARIMAParams, learningRate float64) []float64 {
	numCoeffs := params.P + params.Q
	coeffs := make([]float64, numCoeffs)

	// Initialize coefficients with small random values
	rand.Seed(time.Now().UnixNano())
	for i := range coeffs {
		coeffs[i] = rand.Float64()*0.1 - 0.05
	}

	// Apply differencing
	workingSeries := series
	for d := 0; d < params.D; d++ {
		workingSeries = difference(workingSeries, 1)
	}

	maxLag := max(params.P, params.Q)
	if len(workingSeries) <= maxLag {
		return coeffs
	}

	bestCoeffs := make([]float64, numCoeffs)
	copy(bestCoeffs, coeffs)
	bestError := math.Inf(1)

	// Gradient descent
	for iter := 0; iter < params.MaxIter; iter++ {
		gradients := make([]float64, numCoeffs)
		totalError := 0.0
		numSamples := 0

		for t := maxLag; t < len(workingSeries)-1; t++ {
			prediction := 0.0

			// AR components
			for i := 0; i < params.P; i++ {
				if t-i >= 0 {
					prediction += coeffs[i] * workingSeries[t-i]
				}
			}

			// MA components
			for i := 0; i < params.Q; i++ {
				if t-i >= 0 {
					prediction += coeffs[i+params.P] * workingSeries[t-i]
				}
			}

			// Calculate error
			error := workingSeries[t+1] - prediction
			totalError += error * error
			numSamples++

			// Update gradients
			for i := 0; i < params.P; i++ {
				if t-i >= 0 {
					gradients[i] += -2 * error * workingSeries[t-i]
				}
			}

			for i := 0; i < params.Q; i++ {
				if t-i >= 0 {
					gradients[i+params.P] += -2 * error * workingSeries[t-i]
				}
			}
		}

		// Calculate average error
		if numSamples > 0 {
			totalError /= float64(numSamples)
		}

		// Store best coefficients
		if totalError < bestError {
			bestError = totalError
			copy(bestCoeffs, coeffs)
		}

		// Update coefficients
		for i := range coeffs {
			coeffs[i] -= learningRate * gradients[i] / float64(numSamples)
		}

		// Early stopping
		if totalError < params.Tolerance {
			break
		}
	}

	return bestCoeffs
}
