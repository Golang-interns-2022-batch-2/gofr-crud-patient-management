package patient

import "testing"

func TestValidateName(t *testing.T) {
	tests := []struct {
		desc     string
		name     string
		expected bool
	}{
		{
			desc:     "Case1",
			name:     "aakanksha",
			expected: true,
		},
		{
			desc:     "Case2",
			name:     "",
			expected: false,
		},
	}

	for _, test := range tests {
		test2 := test
		t.Run(test.desc, func(t *testing.T) {
			isValid := validatename(test2.name)
			if isValid != test2.expected {
				t.Errorf("Expected: %v, Got: %v", test2.expected, isValid)
			}
		})
	}
}
