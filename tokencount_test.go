package tokencount

import "testing"

func TestCountTokens(t *testing.T) {
	tests := []struct {
		input string
		min   int
		max   int
	}{
		{"", 0, 0},
		{"hello", 1, 2},
		{"hello world", 2, 4},
		{"The quick brown fox jumps over the lazy dog", 9, 15},
	}
	for _, tt := range tests {
		got := CountTokens(tt.input)
		if got < tt.min || got > tt.max {
			t.Errorf("CountTokens(%q) = %d, want between %d and %d", tt.input, got, tt.min, tt.max)
		}
	}
}

func TestGetModel(t *testing.T) {
	m, ok := GetModel("gpt-4o")
	if !ok {
		t.Fatal("expected gpt-4o to exist")
	}
	if m.MaxTokens != 128000 {
		t.Errorf("expected 128000 max tokens, got %d", m.MaxTokens)
	}

	_, ok = GetModel("nonexistent")
	if ok {
		t.Error("expected nonexistent model to not be found")
	}
}

func TestEstimateCost(t *testing.T) {
	est, err := EstimateCost("gpt-4o", 1000, 500)
	if err != nil {
		t.Fatal(err)
	}
	if est.TotalCost <= 0 {
		t.Error("expected positive cost")
	}

	_, err = EstimateCost("fake-model", 100, 100)
	if err == nil {
		t.Error("expected error for unknown model")
	}
}

func TestFitsContext(t *testing.T) {
	fits, remaining := FitsContext("gpt-4", 4000)
	if !fits {
		t.Error("expected 4000 tokens to fit in gpt-4")
	}
	if remaining != 4192 {
		t.Errorf("expected 4192 remaining, got %d", remaining)
	}

	fits, _ = FitsContext("gpt-4", 10000)
	if fits {
		t.Error("expected 10000 tokens to not fit in gpt-4")
	}
}
