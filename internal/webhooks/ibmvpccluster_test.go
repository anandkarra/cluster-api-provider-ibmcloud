package webhooks_test

import (
	"testing"

	"k8s.io/utils/ptr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/internal/webhooks"
)

func TestValidateIBMVPCCluster_SubnetZoneAndLoadBalancer(t *testing.T) {
	tests := []struct {
		name      string
		cluster   *v1beta2.IBMVPCCluster
		wantError bool
	}{
		{
			name: "Missing zone in control plane subnet",
			cluster: &v1beta2.IBMVPCCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				Spec: v1beta2.IBMVPCClusterSpec{
					Network: &v1beta2.VPCNetworkSpec{
						ControlPlaneSubnets: []v1beta2.Subnet{{ID: nil, Zone: nil}},
						LoadBalancers:       []v1beta2.VPCLoadBalancerSpec{{Name: "lb1"}},
					},
				},
			},
			wantError: true,
		},
		{
			name: "Missing zone in worker subnet",
			cluster: &v1beta2.IBMVPCCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				Spec: v1beta2.IBMVPCClusterSpec{
					Network: &v1beta2.VPCNetworkSpec{
						WorkerSubnets: []v1beta2.Subnet{{ID: nil, Zone: nil}},
						LoadBalancers: []v1beta2.VPCLoadBalancerSpec{{Name: "lb1"}},
					},
				},
			},
			wantError: true,
		},
		{
			name: "No load balancer when network specified",
			cluster: &v1beta2.IBMVPCCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				Spec: v1beta2.IBMVPCClusterSpec{
					Network: &v1beta2.VPCNetworkSpec{
						ControlPlaneSubnets: []v1beta2.Subnet{{ID: ptr.To("id"), Zone: nil}},
						WorkerSubnets:       []v1beta2.Subnet{{ID: ptr.To("id"), Zone: nil}},
						LoadBalancers:       nil,
					},
				},
			},
			wantError: true,
		},
		{
			name: "Valid network spec",
			cluster: &v1beta2.IBMVPCCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				Spec: v1beta2.IBMVPCClusterSpec{
					Network: &v1beta2.VPCNetworkSpec{
						ControlPlaneSubnets: []v1beta2.Subnet{{ID: ptr.To("id"), Zone: nil}},
						WorkerSubnets:       []v1beta2.Subnet{{Zone: ptr.To("zone")}},
						LoadBalancers:       []v1beta2.VPCLoadBalancerSpec{{Name: "lb1"}},
					},
					ControlPlaneLoadBalancer: &v1beta2.VPCLoadBalancerSpec{Name: "lb1"},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := webhooks.ValidateIBMVPCCluster(tt.cluster)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
