package functions

import aimodels "api/ai_models"

// predictARIMA generates predictions using the ARIMA model
func predictARIMA(series []float64, params aimodels.ARIMAParams, coefficients []float64) map[string][]float64 {
    // Store original series for inverse differencing
    originalSeries := make([]float64, len(series))
    copy(originalSeries, series)

    // Apply differencing
    workingSeries := series
    for d := 0; d < params.D; d++ {
        workingSeries = difference(workingSeries, 1)
    }

    maxLag := max(params.P, params.Q)
    if len(workingSeries) <= maxLag {
        return nil
    }

    predictions := make([]float64, len(workingSeries))
    errors := make([]float64, len(workingSeries)) // Store prediction errors for MA terms

    // First maxLag values are set to the actual values
    for i := 0; i < maxLag; i++ {
        predictions[i] = workingSeries[i]
    }

    // Generate predictions
    for t := maxLag; t < len(workingSeries); t++ {
        prediction := 0.0

        // AR components
        for i := 0; i < params.P && t-i-1 >= 0; i++ {
            prediction += coefficients[i] * workingSeries[t-i-1]
        }

        // MA components
        for i := 0; i < params.Q && t-i-1 >= 0; i++ {
            prediction += coefficients[i+params.P] * errors[t-i-1]
        }

        predictions[t] = prediction
        errors[t] = workingSeries[t] - prediction
    }

    // Apply inverse differencing
    result := predictions
    for d := 0; d < params.D; d++ {
        result = inverseDifference(result, originalSeries, 1)
    }

    // Generate extended predictions for 1 day, 1 week, and 1 month
    extendedPredictions := map[string][]float64{
        "1d": extendPredictions(result, coefficients, params, 24),   // 24 hours
        "1w": extendPredictions(result, coefficients, params, 24*7), // 7 days
        "1m": extendPredictions(result, coefficients, params, 24*30), // 30 days
    }

    return extendedPredictions
}

// extendPredictions generates additional predictions for a given period
func extendPredictions(series []float64, coefficients []float64, params aimodels.ARIMAParams, steps int) []float64 {
    extended := make([]float64, len(series)+steps)
    copy(extended, series)

    errors := make([]float64, len(extended)) // Store prediction errors for MA terms

    for t := len(series); t < len(extended); t++ {
        prediction := 0.0

        // AR components
        for i := 0; i < params.P && t-i-1 >= 0; i++ {
            prediction += coefficients[i] * extended[t-i-1]
        }

        // MA components
        for i := 0; i < params.Q && t-i-1 >= 0; i++ {
            prediction += coefficients[i+params.P] * errors[t-i-1]
        }

        extended[t] = prediction
        errors[t] = 0 // No actual data for future errors
    }

    return extended[len(series):]
}