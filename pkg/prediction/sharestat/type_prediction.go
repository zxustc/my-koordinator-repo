package sharestate

import analysisv1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"

type targetPrediction struct {
	// Pod selector that determines which Pods are controlled by this object.
	PodSelector analysisv1.AnalysisTarget
	// From Where to get? One item corresponse to one metric
	MetricSpec []analysisv1.MetricSpec
	// Profilers that need to profile
	Profilers []analysisv1.Profiler
	// Pod/Container dimension
	Level analysisv1.ProfileHierarchyLevel
	// If enable CheckPoint strategy
	UseCheckPoint bool
}

// ProfileID
type ProfileID struct {
	// metricPredication CR’s namespace
	Namespace string
	// metricPredication CR’s name
	Name string
}

type ProfilerResult struct {
}

type TargetPrediction interface {
	// fron-end can use channel to solve more than one status updated
	updatePredictionStatus(id ProfileID, result analysisv1.ProfileResult)
}

func (*targetPrediction) updatePredictionStatus(id ProfileID, result analysisv1.ProfileResult) {
	// TODO
	panic("TODO")
}
