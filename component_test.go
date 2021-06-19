package model

import "testing"

func TestStepComponentType(t *testing.T) {
	s := inittests(t)
	step := s.Steps[0]
	if step.Component.Type != "executable" {
		t.Errorf("The step's component type was '%s' when it should have been 'executable'", step.Component.Type)
	}
}

func TestStepComponentName(t *testing.T) {
	s := inittests(t)
	step := s.Steps[0]
	if step.Component.Name != "QATestTool.sh" {
		t.Errorf("The step's component name was '%s' when it should have been 'QATestTool.sh'", step.Component.Name)
	}
}

func TestStepComponentLocation(t *testing.T) {
	s := inittests(t)
	step := s.Steps[0]
	if step.Component.Location != "/usr/local2/bin" {
		t.Errorf("The step's component location was '%s' when it should have been '/usr/local2/bin'", step.Component.Location)
	}
}

func TestStepComponentDescription(t *testing.T) {
	s := inittests(t)
	step := s.Steps[0]
	if step.Component.Description != "Test script to emulate a tool installed" {
		t.Errorf("The step's component description was '%s' when it should have been 'Test script to emulate a tool installed'", step.Component.Description)
	}
}
