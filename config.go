package httptracer

type Config struct {
	ServiceName    string
	ServiceVersion string

	// AnalyticsRate sample rate, default: 1
	//
	// Set as -1 to disable otherwise provide value between 0 and 1
	AnalyticsRate float64
}
