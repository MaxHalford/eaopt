package presets

import "github.com/MaxHalford/gago"

// Alignment returns a configuration for solving string alignment problems.
// function. The output will be a genome of a certain length with genes
// belonging to a corpus of elements.
func Alignment(length int, corpus []string, distance func([]string) float64) gago.GA {
    return gago.GA{
        NbPopulations: 2,
        NbIndividuals: 30,
        NbGenes:       length,
        NbParents:     6,
        Ff: gago.StringFunction{
            Image: distance,
        },
        Initializer: gago.ISUniform{
            Corpus: corpus,
        },
        Selector: gago.STournament{
            NbParticipants: 3,
        },
        Crossover: gago.CPoint{
            NbPoints: 2,
        },
        Mutator: gago.MutPermute{
            Max: 3,
        },
        MutRate:      0.5,
        Migrator:     gago.MigShuffle{},
        MigFrequency: 10,
    }
}
