package resources

import (
	"testing"
)

func Test_parametersEqual(t *testing.T) {
	tests := []struct {
		name     string
		existing map[string]string
		next     map[string]string
		want     bool
	}{
		{
			name:     "existing nil, current non-nil",
			existing: nil,
			next:     map[string]string{"key": "value"},
			want:     false,
		},
		{
			name:     "existing non-nil, current nil",
			existing: nil,
			next:     map[string]string{"key": "value"},
			want:     false,
		},
		{
			name:     "existing empty, current non-empty",
			existing: map[string]string{},
			next:     map[string]string{"key": "value"},
			want:     false,
		},
		{
			name:     "existing non-empty, current empty",
			existing: map[string]string{"key": "value"},
			next:     map[string]string{},
			want:     false,
		},
		{
			name:     "existing and current both empty",
			existing: map[string]string{},
			next:     map[string]string{},
			want:     true,
		},
		{
			name:     "existing and current both non-empty and equal",
			existing: map[string]string{"key": "value"},
			next:     map[string]string{"key": "value"},
			want:     true,
		},
		{
			name:     "existing and current both non-empty and not equal",
			existing: map[string]string{"key": "value"},
			next:     map[string]string{"key": "value2"},
			want:     false,
		},
		{
			name:     "existing and current both non-empty and not equal in length",
			existing: map[string]string{"key": "value"},
			next:     map[string]string{"key": "value", "key2": "value2"},
			want:     false,
		},
		{
			name:     "existing and current both non-empty and not equal in length",
			existing: map[string]string{"key": "value", "key2": "value2"},
			next:     map[string]string{"key": "value"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parametersEqual(tt.existing, tt.next); got != tt.want {
				t.Errorf("parametersEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
