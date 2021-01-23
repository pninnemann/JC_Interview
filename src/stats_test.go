package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestSetup(t *testing.T) {
	totals = make(map[string]ActionTotal)
	input1 := "{\"action\":\"jump\", \"time\":100}"
	input2 := "{\"action\":\"run\", \"time\":75}"
	input3 := "{\"action\":\"jump\", \"time\":200}"

	err := addAction(input1)
	if err == nil {
		err = addAction(input2)
	}
	if err == nil {
		err = addAction(input3)
	}
	if err != nil {
		t.Error("Expected no errors during add Action; Got error ", err)
	}
}

func TestGetStats(t *testing.T) {
	TestSetup(t)

	expected1 := "[{\"action\":\"jump\",\"avg\":150},{\"action\":\"run\",\"avg\":75}]"
	expected2 := "[{\"action\":\"run\",\"avg\":75},{\"action\":\"jump\",\"avg\":150}]"

	ret := getStats()
	if ret != expected1 && ret != expected2 {
		t.Error("getStats did not return expected output")
	}
}

func TestBadInput(t *testing.T)  {
	input1 := "{\"action\":\"jump\", \"time\":onehundred}"
	input2 := "{\"action\":jump, \"time\":100}"
	input3 := "{\"action\":\"jump\" \"time\":200}"

	err := addAction(input1)
	if err == nil {
		t.Error("Expected error thrown")
	}
	err = addAction(input2)
	if err == nil {
		t.Error("Expected error thrown")
	}
	err = addAction(input3)
	if err == nil {
		t.Error("Expected error thrown")
	}
}

func TestNegativeTime(t *testing.T) {
	totals = make(map[string]ActionTotal)
	input1 := "{\"action\":\"jump\", \"time\":-100}"

	err := addAction(input1)
	if err != nil {
		t.Error("Expected no errors during add Action; Got error ", err)
	}

	expected1 := "[{\"action\":\"jump\",\"avg\":-100}]"
	ret := getStats()
	if ret != expected1 {
		t.Error("getStats did not return expected output")
	}

	input2 := "{\"action\":\"jump\", \"time\":100}"

	err = addAction(input2)
	if err != nil {
		t.Error("Expected no errors during add Action; Got error ", err)
	}

	expected2 := "[{\"action\":\"jump\",\"avg\":0}]"
	ret = getStats()
	if ret != expected2 {
		t.Error("getStats did not return expected output")
	}
}

func TestFractionalInput(t *testing.T) {
	totals = make(map[string]ActionTotal)
	input1 := "{\"action\":\"jump\", \"time\":1.5}"
	input2 := "{\"action\":\"jump\", \"time\":2.5}"

	err := addAction(input1)
	if err != nil {
		t.Error("Expected no errors during add Action; Got error ", err)
	}
	err = addAction(input2)
	if err != nil {
		t.Error("Expected no errors during add Action; Got error ", err)
	}

	expected1 := "[{\"action\":\"jump\",\"avg\":2}]"
	ret := getStats()
	if ret != expected1 {
		t.Error("getStats did not return expected output")
	}
}

func TestFractionalAverage(t *testing.T) {
	totals = make(map[string]ActionTotal)
	input1 := "{\"action\":\"jump\", \"time\":1}"
	input2 := "{\"action\":\"jump\", \"time\":2}"

	err := addAction(input1)
	if err != nil {
		t.Error("Expected no errors during add Action; Got error ", err)
	}
	err = addAction(input2)
	if err != nil {
		t.Error("Expected no errors during add Action; Got error ", err)
	}

	expected1 := "[{\"action\":\"jump\",\"avg\":1.5}]"
	ret := getStats()
	if ret != expected1 {
		t.Error("getStats did not return expected output")
	}
}

func TestParallelCalls(t *testing.T) {
	totals = make(map[string]ActionTotal)

	var waitGroup sync.WaitGroup
	waitGroup.Add(3)
	OneHundredjumps := func(t *testing.T) {
		defer waitGroup.Done()
		for i := 1; i <= 100; i++ {
			input := fmt.Sprintf("{\"action\":\"jump\", \"time\":%d}",i)
			err := addAction(input)
			if err != nil {
				t.Error("Expected no errors during add Action; Got error ", err)
				break
			}
		}
	}
	OneHundredruns := func(t *testing.T) {
		defer waitGroup.Done()
		input := "{\"action\":\"run\", \"time\":10}"
		for i := 0; i < 100; i++ {
			err := addAction(input)
			if err != nil {
				t.Error("Expected no errors during add Action; Got error ", err)
				break
			}
		}
	}
	OneHundredskips := func(t *testing.T) {
		defer waitGroup.Done()
		input := "{\"action\":\"skip\", \"time\":500}"
		for i := 0; i < 100; i++ {
			err := addAction(input)
			if err != nil {
				t.Error("Expected no errors during add Action; Got error ", err)
				break
			}
		}
	}

	go OneHundredjumps(t)
	go OneHundredruns(t)
	go OneHundredskips(t)
	waitGroup.Wait()

	expected1 := "[{\"action\":\"skip\",\"avg\":500},{\"action\":\"jump\",\"avg\":50.5},{\"action\":\"run\",\"avg\":10}]"
	expected2 := "[{\"action\":\"skip\",\"avg\":500},{\"action\":\"run\",\"avg\":10},{\"action\":\"jump\",\"avg\":50.5}]"
	expected3 := "[{\"action\":\"jump\",\"avg\":50.5},{\"action\":\"run\",\"avg\":10},{\"action\":\"skip\",\"avg\":500},]"
	expected4 := "[{\"action\":\"jump\",\"avg\":50.5},{\"action\":\"skip\",\"avg\":500},{\"action\":\"run\",\"avg\":10}]"
	expected5 := "[{\"action\":\"run\",\"avg\":10},{\"action\":\"skip\",\"avg\":500},{\"action\":\"jump\",\"avg\":50.5}]"
	expected6 := "[{\"action\":\"run\",\"avg\":10},{\"action\":\"jump\",\"avg\":50.5},{\"action\":\"skip\",\"avg\":500}]"

	ret := getStats()
	if ret != expected1 && ret != expected2 && ret != expected3 && ret != expected4 && ret != expected5 && ret != expected6 {
		t.Error("getStats did not return expected output")
		log.Println(ret)
	}
}