package eaopt

import (
	"fmt"
	"testing"
)

func TestNewGARNGNotNil(t *testing.T) {
	var conf = NewDefaultGAConfig()
	conf.RNG = nil
	var ga, err = conf.NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if ga.RNG == nil {
		t.Error("RNG should not be nil")
	}
}

func TestNewGAErrors(t *testing.T) {
	var testCases = []struct {
		conf GAConfig
	}{
		{func() GAConfig { c := NewDefaultGAConfig(); c.NPops = 0; return c }()},
		{func() GAConfig { c := NewDefaultGAConfig(); c.PopSize = 0; return c }()},
		{func() GAConfig { c := NewDefaultGAConfig(); c.NGenerations = 0; return c }()},
		{func() GAConfig { c := NewDefaultGAConfig(); c.HofSize = 0; return c }()},
		{func() GAConfig { c := NewDefaultGAConfig(); c.Model = nil; return c }()},
		{func() GAConfig { c := NewDefaultGAConfig(); c.Model = ModValidateError{}; return c }()},
		{func() GAConfig { c := NewDefaultGAConfig(); c.Migrator = MigRing{0}; return c }()},
		{func() GAConfig { c := NewDefaultGAConfig(); c.Migrator = MigRing{1}; c.MigFrequency = 0; return c }()},
		{func() GAConfig { c := NewDefaultGAConfig(); c.Speciator = SpecValidateError{}; return c }()},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var _, err = tc.conf.NewGA()
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
	}
}

func TestNewDefaultGAConfig(t *testing.T) {
	var _, err = NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
