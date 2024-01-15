package ecosystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestBlueprintSpec_CheckDoguHealth(t *testing.T) {
//	tests := []struct {
//		name               string
//		inputSpec          *BlueprintSpec
//		installedDogus     map[string]*ecosystem.DoguInstallation
//		expectedStatus     StatusPhase
//		expectedEventNames []string
//		expectedEventMsgs  []string
//	}{
//		{
//			name:               "should return early if ignore dogu health is configured",
//			inputSpec:          &BlueprintSpec{Config: BlueprintConfiguration{IgnoreDoguHealth: true}},
//			installedDogus:     nil,
//			expectedStatus:     StatusPhaseIgnoreDoguHealth,
//			expectedEventNames: []string{"IgnoreDoguHealth"},
//			expectedEventMsgs:  []string{"ignore dogu health flag is set; ignoring dogu health"},
//		},
//		{
//			name:      "should write unhealthy dogus in event",
//			inputSpec: &BlueprintSpec{},
//			installedDogus: map[string]*ecosystem.DoguInstallation{
//				"ldap": {
//					Namespace: "official",
//					Name:      "ldap",
//					Version:   mustParseVersion("1.2.3-1"),
//					Health:    "unavailable",
//				},
//				"postfix": {
//					Namespace: "official",
//					Name:      "postfix",
//					Version:   mustParseVersion("5.6.7-5"),
//					Health:    "available",
//				},
//				"nginx-ingress": {
//					Namespace: "k8s",
//					Name:      "nginx-ingress",
//					Version:   mustParseVersion("2.3.4-2"),
//					Health:    "unknownError",
//				},
//				"nginx-static": {
//					Namespace: "k8s",
//					Name:      "nginx-static",
//					Version:   mustParseVersion("2.3.4-2"),
//					Health:    "available",
//				},
//			},
//			expectedStatus:     StatusPhaseEcosystemUnhealthyUpfront,
//			expectedEventNames: []string{"DogusUnhealthy"},
//			expectedEventMsgs:  []string{"2 dogus are unhealthy: k8s/nginx-ingress:2.3.4-2 is unknownError, official/ldap:1.2.3-1 is unavailable"},
//		},
//		{
//			name:      "all dogus healthy",
//			inputSpec: &BlueprintSpec{},
//			installedDogus: map[string]*ecosystem.DoguInstallation{
//				"ldap": {
//					Namespace: "official",
//					Name:      "ldap",
//					Version:   mustParseVersion("1.2.3-1"),
//					Health:    "available",
//				},
//				"postfix": {
//					Namespace: "official",
//					Name:      "postfix",
//					Version:   mustParseVersion("5.6.7-5"),
//					Health:    "available",
//				},
//				"nginx-ingress": {
//					Namespace: "k8s",
//					Name:      "nginx-ingress",
//					Version:   mustParseVersion("2.3.4-2"),
//					Health:    "available",
//				},
//				"nginx-static": {
//					Namespace: "k8s",
//					Name:      "nginx-static",
//					Version:   mustParseVersion("2.3.4-2"),
//					Health:    "available",
//				},
//			},
//			expectedStatus:     StatusPhaseEcosystemHealthyUpfront,
//			expectedEventNames: []string{"DogusHealthy"},
//			expectedEventMsgs:  []string{""},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.inputSpec.CheckEcosystemHealthUpfront(tt.installedDogus)
//			eventNames := util.Map(tt.inputSpec.Events, Event.Name)
//			eventMsgs := util.Map(tt.inputSpec.Events, Event.Message)
//			assert.ElementsMatch(t, tt.expectedEventNames, eventNames)
//			assert.ElementsMatch(t, tt.expectedEventMsgs, eventMsgs)
//		})
//	}
//}

func TestCalculateDoguHealthResult(t *testing.T) {
	tests := []struct {
		name  string
		dogus []*DoguInstallation
		want  DoguHealthResult
	}{
		{
			name: "",
			dogus: []*DoguInstallation{
				{
					Name:   "postgresql",
					Health: AvailableHealthStatus,
				},
				{
					Name:   "postfix",
					Health: UnavailableHealthStatus,
				},
				{
					Name:   "ldap",
					Health: PendingHealthStatus,
				},
			},
			want: DoguHealthResult{
				DogusByStatus: map[HealthStatus][]DoguName{
					AvailableHealthStatus:   {"postgresql"},
					UnavailableHealthStatus: {"postfix"},
					PendingHealthStatus:     {"ldap"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CalculateDoguHealthResult(tt.dogus), "CalculateDoguHealthResult(%v)", tt.dogus)
		})
	}
}
