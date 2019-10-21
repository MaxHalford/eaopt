package eaopt

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"time"
)

// GAConfig contains fields that are necessary to instantiate a GA.
type GAConfig struct {
	// Required fields
	NPops        uint
	PopSize      uint
	NGenerations uint
	HofSize      uint
	Model        Model

	// Optional fields
	ParallelEval bool // Whether to evaluate Individuals in parallel or not
	Migrator     Migrator
	MigFrequency uint // Frequency at which migrations occur
	Speciator    Speciator
	Logger       *log.Logger
	Callback     func(ga *GA)
	EarlyStop    func(ga *GA) bool
	RNG          *rand.Rand

	// Optional, marshaling fields
	PopulationsReader     io.Reader
	GenomeJSONUnmarshaler func([]byte) (Genome, error)
}

// NewGA returns a pointer to a GA instance and checks for configuration
// errors.
func (conf GAConfig) NewGA() (*GA, error) {
	// Check for default values
	if conf.RNG == nil {
		conf.RNG = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	// Check the configuration is valid
	if conf.NPops == 0 {
		return nil, errors.New("NPops has to be strictly higher than 0")
	}
	if conf.PopSize == 0 {
		return nil, errors.New("PopSize has to be strictly higher than 0")
	}
	if conf.NGenerations == 0 {
		return nil, errors.New("NGenerations has to be strictly higher than 0")
	}
	if conf.HofSize == 0 {
		return nil, errors.New("HofSize has to be strictly higher than 0")
	}
	if conf.Model == nil {
		return nil, errors.New("Model has to be provided")
	}
	if modelErr := conf.Model.Validate(); modelErr != nil {
		return nil, modelErr
	}
	if conf.Migrator != nil {
		if migErr := conf.Migrator.Validate(); migErr != nil {
			return nil, migErr
		}
		if conf.MigFrequency == 0 {
			return nil, errors.New("MigFrequency should be higher than 0")
		}
	}
	if conf.Speciator != nil {
		if specErr := conf.Speciator.Validate(); specErr != nil {
			return nil, specErr
		}
	}
	// Initialize the GA
	return &GA{GAConfig: conf}, nil
}

// NewDefaultGAConfig returns a valid GAConfig with default values.
func NewDefaultGAConfig() GAConfig {
	return GAConfig{
		NPops:        1,
		PopSize:      30,
		HofSize:      1,
		NGenerations: 50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate:   0.5,
			CrossRate: 0.7,
		},
		ParallelEval: false,
	}
}
