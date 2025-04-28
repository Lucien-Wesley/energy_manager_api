package functions

import (
	aimodels "api/ai_models"
	"fmt"
	"math"
)

func Predictions(trainData []float64) map[string][]float64 {

	// Normalize training data
	trainMean := 0.0
	for _, v := range trainData {
		trainMean += v
	}
	trainMean /= float64(len(trainData))

	trainStd := 0.0
	for _, v := range trainData {
		trainStd += (v - trainMean) * (v - trainMean)
	}
	trainStd = math.Sqrt(trainStd / float64(len(trainData)))

	normalizedTrain := make([]float64, len(trainData))
	for i, v := range trainData {
		normalizedTrain[i] = (v - trainMean) / trainStd
	}

	// ARIMA parameters
	params := aimodels.ARIMAParams{
		P:         2, // AR order
		D:         1, // Differencing order
		Q:         2, // MA order
		Tolerance: 1e-6,
		MaxIter:   1000,
	}

	// Train model
	learningRate := 0.01
	optimizedCoefficients := optimizeARIMA(normalizedTrain, params, learningRate)

	// Generate predictions using the provided predictARIMA function
	predictionResults := predictARIMA(normalizedTrain, params, optimizedCoefficients)
	if predictionResults == nil {
		fmt.Println("Error: Failed to generate predictions")
		return nil
	}

	// Denormalize predictions
	denormalizedResults := make(map[string][]float64)
	for key, preds := range predictionResults {
		denormPreds := make([]float64, len(preds))
		for i := range preds {
			denormPreds[i] = preds[i]*trainStd + trainMean
		}
		denormalizedResults[key] = denormPreds
	}

	// Calculate and print training metrics for the training data predictions
	trainPredictions := denormalizedResults["1d"] // Assuming "1d" corresponds to training predictions
	trainMAE, trainMSE, trainRMSE := calculateMetrics(trainData, trainPredictions)
	fmt.Printf("\nTraining Metrics:\n")
	fmt.Printf("MAE: %.4f\n", trainMAE)
	fmt.Printf("MSE: %.4f\n", trainMSE)
	fmt.Printf("RMSE: %.4f\n", trainRMSE)

	return denormalizedResults
}
