package resources

import (
	"testing"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func Test_firewallsSame(t *testing.T) {
	tests := []struct {
		name string
		ls1  []ccx.FirewallRule
		ls2  []ccx.FirewallRule
		want bool
	}{
		{
			name: "both empty",
			ls1:  nil,
			ls2:  nil,
			want: true,
		},
		{
			name: "ls1 empty",
			ls1:  nil,
			ls2: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
			},
			want: false,
		},
		{
			name: "ls2 empty",
			ls1: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
			},
			ls2:  nil,
			want: false,
		},
		{
			name: "ls1 subset of ls2",
			ls1: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
			},
			ls2: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
				{Source: "1.2.3.5/32", Description: "bar"},
			},
			want: false,
		},
		{
			name: "ls2 subset of ls1",
			ls1: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
				{Source: "1.2.3.5/32", Description: "bar"},
			},
			ls2: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
			},
			want: false,
		},
		{
			name: "ls1 and ls2 have different elements",
			ls1: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
			},
			ls2: []ccx.FirewallRule{
				{Source: "1.2.3.5/32", Description: "bar"},
			},
			want: false,
		},
		{
			name: "ls1 and ls2 have same elements, different order",
			ls1: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
				{Source: "1.2.3.5/32", Description: "bar"},
			},
			ls2: []ccx.FirewallRule{
				{Source: "1.2.3.5/32", Description: "bar"},
				{Source: "1.2.3.4/32", Description: "foo"},
			},
			want: true,
		},
		{
			name: "ls1 and ls2 have same elements, same order",
			ls1: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
				{Source: "1.2.3.5/32", Description: "bar"},
			},
			ls2: []ccx.FirewallRule{
				{Source: "1.2.3.4/32", Description: "foo"},
				{Source: "1.2.3.5/32", Description: "bar"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firewallsSame(tt.ls1, tt.ls2); got != tt.want {
				t.Errorf("firewallsChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}
