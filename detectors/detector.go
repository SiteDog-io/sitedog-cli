package detectors

// Detector interface for project analysis
type Detector interface {
	Name() string
	Description() string
	Detect() ([]*DetectionResult, error)
	ShouldRun() bool
}

// DetectionResult represents what a detector found
type DetectionResult struct {
	Key         string
	Value       interface{}
	Description string
	Confidence  float64 // 0.0 to 1.0
}