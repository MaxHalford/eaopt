package eaopt

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"math/rand"
	"reflect"
	"testing"
)

func TestPopLog(t *testing.T) {
	var (
		pop    = newPopulation(42, NewVector, rand.New(rand.NewSource(42)))
		b      bytes.Buffer
		logger = log.New(&b, "", 0)
	)
	pop.Individuals.Evaluate(false)
	pop.Log(logger)
	var expected = "pop_id=KVm min=-21.342844 max=18.440761 avg=-1.404246 std=11.739691\n"
	if s := b.String(); s != expected {
		t.Errorf("Expected %s, got %s", expected, s)
	}
}

func TestPopJSONMarshal(t *testing.T) {
	pop1 := newPopulation(42, NewVector, rand.New(rand.NewSource(42)))
	pop1.Individuals.Evaluate(false)
	encodedPop1, err := json.Marshal(pop1)
	if err != nil {
		t.Fatal(err)
	}

	var pop2 Population
	pop2.JSONUnmarshaler = VectorJSONUnmarshaler
	err = json.Unmarshal(encodedPop1, &pop2)
	if err != nil {
		t.Fatal(err)
	}
	encodedPop2, err := json.Marshal(pop2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(encodedPop1, encodedPop2) {
		t.Fatal("Marshaling error")
	}
}

func TestPopGOBMarshal(t *testing.T) {
	pop1 := newPopulation(42, NewVector, rand.New(rand.NewSource(42)))
	pop1.Individuals.Evaluate(false)

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&pop1); err != nil {
		t.Fatal(err)
	}

	var pop2 Population
	decoder := gob.NewDecoder(&buf)
	err := decoder.Decode(&pop2)
	if err != nil {
		t.Fatal(err)
	}
	pop2.Individuals.Evaluate(false)
	if pop1.String() != pop2.String() {
		t.Errorf("Expected %s, got %s", pop1.String(), pop2.String())
	}
}

func TestPopsJSONMarshal(t *testing.T) {
	pop1 := newPopulation(3, NewVector, rand.New(rand.NewSource(42)))
	pop1.Individuals.Evaluate(false)
	pop2 := newPopulation(3, NewVector, rand.New(rand.NewSource(201)))
	pop2.Individuals.Evaluate(false)

	pops := Populations{pop1, pop2}
	encodedPops, err := json.Marshal(pops)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer(encodedPops)
	decodedPops, err := newPopulationsFromReader(uint(len(pops)), buf, rand.New(rand.NewSource(42)), VectorJSONUnmarshaler)
	if err != nil {
		t.Fatal(err)
	}

	encodedDecodedPops, err := json.Marshal(decodedPops)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(encodedPops, encodedDecodedPops) {
		t.Fatal("Marshaling error")
	}
}

func TestPopsGOBMarshal(t *testing.T) {
	pop1 := newPopulation(3, NewVector, rand.New(rand.NewSource(42)))
	pop1.Individuals.Evaluate(false)
	pop2 := newPopulation(3, NewVector, rand.New(rand.NewSource(201)))
	pop2.Individuals.Evaluate(false)

	pops := Populations{pop1, pop2}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&pops); err != nil {
		t.Fatal(err)
	}

	decodedPops, err := newPopulationsFromReader(uint(len(pops)), &buf, rand.New(rand.NewSource(42)), nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := range pops {
		if !reflect.DeepEqual(pops[i].Individuals, decodedPops[i].Individuals) {
			t.Fatal("Marshaling error")
		}
	}
}
