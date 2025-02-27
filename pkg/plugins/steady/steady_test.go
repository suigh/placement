package steady

import (
	"context"
	"testing"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapiv1 "open-cluster-management.io/api/cluster/v1"
	clusterapiv1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
	testinghelpers "open-cluster-management.io/placement/pkg/helpers/testing"
)

func TestScoreClusterWithSteady(t *testing.T) {
	cases := []struct {
		name              string
		placement         *clusterapiv1alpha1.Placement
		clusters          []*clusterapiv1.ManagedCluster
		existingDecisions []runtime.Object
		expectedScores    map[string]int64
	}{
		{
			name:      "no decisions",
			placement: testinghelpers.NewPlacement("test", "test").Build(),
			clusters: []*clusterapiv1.ManagedCluster{
				testinghelpers.NewManagedCluster("cluster1").Build(),
				testinghelpers.NewManagedCluster("cluster2").Build(),
				testinghelpers.NewManagedCluster("cluster3").Build(),
			},
			existingDecisions: []runtime.Object{},
			expectedScores:    map[string]int64{"cluster1": 0, "cluster2": 0, "cluster3": 0},
		},
		{
			name:      "one decisions",
			placement: testinghelpers.NewPlacement("test", "test").Build(),
			clusters: []*clusterapiv1.ManagedCluster{
				testinghelpers.NewManagedCluster("cluster1").Build(),
				testinghelpers.NewManagedCluster("cluster2").Build(),
				testinghelpers.NewManagedCluster("cluster3").Build(),
			},
			existingDecisions: []runtime.Object{
				testinghelpers.NewPlacementDecision("test", "test1").WithLabel(placementLabel, "test").WithDecisions("cluster1").Build(),
			},
			expectedScores: map[string]int64{"cluster1": 100, "cluster2": 0, "cluster3": 0},
		},
		{
			name:      "one decisions",
			placement: testinghelpers.NewPlacement("test", "test").Build(),
			clusters: []*clusterapiv1.ManagedCluster{
				testinghelpers.NewManagedCluster("cluster1").Build(),
				testinghelpers.NewManagedCluster("cluster2").Build(),
				testinghelpers.NewManagedCluster("cluster3").Build(),
			},
			existingDecisions: []runtime.Object{
				testinghelpers.NewPlacementDecision("test", "test1").WithLabel(placementLabel, "test").WithDecisions("cluster1").Build(),
				testinghelpers.NewPlacementDecision("test", "test2").WithLabel(placementLabel, "test").WithDecisions("cluster3").Build(),
			},
			expectedScores: map[string]int64{"cluster1": 100, "cluster2": 0, "cluster3": 100},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			steady := &Steady{
				handle: testinghelpers.NewFakePluginHandle(t, nil, c.existingDecisions...),
			}

			scores, err := steady.Score(context.TODO(), c.placement, c.clusters)
			if err != nil {
				t.Errorf("Expect no error, but got %v", err)
			}

			if !apiequality.Semantic.DeepEqual(scores, c.expectedScores) {
				t.Errorf("Expect score %v, but got %v", c.expectedScores, scores)
			}
		})
	}
}
