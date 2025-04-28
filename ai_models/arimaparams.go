package aimodels

// ARIMAParams holds the parameters for the ARIMA model
type ARIMAParams struct {
	P         int     // Autoregressive order
	D         int     // Differencing order
	Q         int     // Moving average order
	Tolerance float64 // Convergence tolerance
	MaxIter   int     // Maximum number of iterations
}
