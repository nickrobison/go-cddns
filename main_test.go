package main

import "testing"

func TestConfigValidation(t *testing.T) {
	// Test empty key
	err := validateConfig(&config{
		UpdateInterval: 15,
		Key:            "",
		Email:          "hello",
	})
	if err == nil {
		t.Errorf("API Key cannot be empty")
	}
	if err.Error() != "Empty API Key, exiting" {
		t.Errorf("Wrong error message")
	}

	err = validateConfig(&config{
		UpdateInterval: 15,
		Email:          "hello",
	})
	if err == nil {
		t.Errorf("API Key cannot be nil")
	}
	if err.Error() != "Empty API Key, exiting" {
		t.Errorf("Wrong error message")
	}

	// Test empty email
	err = validateConfig(&config{
		UpdateInterval: 15,
		Key:            "test-key",
		Email:          "",
	})
	if err == nil {
		t.Errorf("Email cannot be empty")
	}
	if err.Error() != "Empty Email Address, exiting" {
		t.Errorf("Wrong error message")
	}

	err = validateConfig(&config{
		UpdateInterval: 15,
		Key:            "test-key",
	})
	if err == nil {
		t.Errorf("Email address cannot be null")
	}
	if err.Error() != "Empty Email Address, exiting" {
		t.Errorf("Wrong error message")
	}

	// Test domain name
	err = validateConfig(&config{
		UpdateInterval: 15,
		Key:            "test-key",
		Email:          "test-email",
		DomainName:     "",
	})
	if err == nil {
		t.Errorf("Email cannot be empty")
	}
	if err.Error() != "Empty Domain Name, exiting" {
		t.Errorf("Wrong error message")
	}

	err = validateConfig(&config{
		UpdateInterval: 15,
		Key:            "test-key",
		Email:          "test-email",
	})
	if err == nil {
		t.Errorf("Domain name cannot be null")
	}
	if err.Error() != "Empty Domain Name, exiting" {
		t.Errorf("Wrong error message")
	}

	// Test records
	err = validateConfig(&config{
		UpdateInterval: 15,
		Key:            "test-key",
		Email:          "test-email",
		DomainName:     "test-domain",
		RecordNames:    []string{},
	})
	if err == nil {
		t.Errorf("Records cannot be empty")
	}
	if err.Error() != "No records to update, exiting" {
		t.Errorf("Wrong error message")
	}

	err = validateConfig(&config{
		UpdateInterval: 15,
		Key:            "test-key",
		Email:          "test-email",
		DomainName:     "test-domain",
	})
	if err == nil {
		t.Errorf("Records cannot be nil")
	}
	if err.Error() != "No records to update, exiting" {
		t.Errorf("Wrong error message")
	}

	// Update interval
	err = validateConfig(&config{
		UpdateInterval: 3,
		Key:            "test-key",
		Email:          "test-email",
		DomainName:     "test-domain",
		RecordNames:    []string{"test-record"},
	})
	if err == nil {
		t.Errorf("Update interval should fail")
	}
	if err.Error() != "Update interval cannot be less than 5 minutes" {
		t.Errorf("Wrong error message")
	}

	err = validateConfig(&config{
		Key:         "test-key",
		Email:       "test-email",
		DomainName:  "test-domain",
		RecordNames: []string{"test-record"},
	})
	if err == nil {
		t.Errorf("Update interval should fail")
	}
}
