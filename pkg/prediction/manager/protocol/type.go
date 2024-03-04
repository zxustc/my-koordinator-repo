package protocol

import (
	analysisv1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ PredictionProfileKey = &predictionProfileKey{}

// predictionProfileKey is used to collect infos to position a unique profile task
type predictionProfileKey struct {
	TargetType analysisv1.AnalysisTargetType
	//Pod selector that determines profile pods, used with type "workload".
	Workload *analysisv1.WorkloadRef
	//Pod selector that determines profile pods, used with type "podselector".
	PodSelector *analysisv1.PodSelectorRef
	//From Where to get? One item corresponse to one metric
	MetricSpec analysisv1.MetricSpec
	//Profilers that need to profile
	Profilers []analysisv1.Profiler
	//Profile Hierarchy
	Stage analysisv1.ProfileHierarchy
	//If enable CheckPoint strategy
	UseCheckPoint bool
}

// ProfileID
type ProfileID struct {
	//metricPredication CR’s namespace
	Namespace string
	//metricPredication CR’s name
	Name string
}

func (p *predictionProfileKey) Type() string {
	return string(p.TargetType)
}

func (p *predictionProfileKey) StageName() string {
	return string(p.Stage.Level)
}

func (p *predictionProfileKey) WorkloadName() string {
	if p.TargetType == "workload" {
		return p.Workload.Kind
	} else {
		return ""
	}
}

func (p *predictionProfileKey) WorkloadKind() string {
	if p.TargetType == "workload" {
		return p.Workload.Name
	} else {
		return ""
	}
}

func (p *predictionProfileKey) LabelSelector() *metav1.LabelSelector {
	if p.TargetType == "workload" {
		return nil
	} else {
		return p.PodSelector.Selector
	}
}

func (p *predictionProfileKey) MetricSourceName() string {
	return string(p.MetricSpec.Source)
}

func (p *predictionProfileKey) ProfilerName() []string {
	var profilers []string
	for _, Profiler := range p.Profilers {
		profilers = append(profilers, Profiler.Name)
	}
	return profilers
}

func InitPredictionProfileKey(metricPrediction analysisv1.MetricPredictionSpec) *predictionProfileKey {
	predictionProfileKey := &predictionProfileKey{
		TargetType:  metricPrediction.Target.Type,
		Workload:    metricPrediction.Target.Workload,
		PodSelector: metricPrediction.Target.PodSelector,
		MetricSpec:  metricPrediction.Metric,
		Profilers:   metricPrediction.Profilers,
	}
	return predictionProfileKey
}

type PredictionProfileState struct {
}
