package main

import "testing"

func TestParseRelaySendArgs_ParsesTimeoutOverride(t *testing.T) {
	t.Setenv("CC_PROJECT", "source")
	t.Setenv("CC_SESSION_KEY", "telegram:chat:user")

	opts, err := parseRelaySendArgs([]string{
		"--to", "gemini",
		"--timeout", "300",
		"please",
		"respond",
	})
	if err != nil {
		t.Fatalf("parseRelaySendArgs returned error: %v", err)
	}
	if opts.from != "source" {
		t.Fatalf("opts.from = %q, want %q", opts.from, "source")
	}
	if opts.to != "gemini" {
		t.Fatalf("opts.to = %q, want %q", opts.to, "gemini")
	}
	if opts.message != "please respond" {
		t.Fatalf("opts.message = %q, want %q", opts.message, "please respond")
	}
	if opts.timeoutSecs == nil || *opts.timeoutSecs != 300 {
		t.Fatalf("opts.timeoutSecs = %v, want 300", opts.timeoutSecs)
	}
}

func TestParseRelaySendArgs_ParsesTimeoutAliasAndZero(t *testing.T) {
	t.Setenv("CC_PROJECT", "source")
	t.Setenv("CC_SESSION_KEY", "telegram:chat:user")

	opts, err := parseRelaySendArgs([]string{
		"--to", "gemini",
		"--timeout-secs", "0",
		"--message", "hello",
	})
	if err != nil {
		t.Fatalf("parseRelaySendArgs returned error: %v", err)
	}
	if opts.timeoutSecs == nil || *opts.timeoutSecs != 0 {
		t.Fatalf("opts.timeoutSecs = %v, want 0", opts.timeoutSecs)
	}
}

func TestParseRelaySendArgs_RejectsNegativeTimeout(t *testing.T) {
	t.Setenv("CC_PROJECT", "source")
	t.Setenv("CC_SESSION_KEY", "telegram:chat:user")

	_, err := parseRelaySendArgs([]string{
		"--to", "gemini",
		"--timeout", "-1",
		"--message", "hello",
	})
	if err == nil {
		t.Fatal("expected error for negative timeout, got nil")
	}
	if err.Error() != "timeout must be >= 0" {
		t.Fatalf("err = %q, want %q", err.Error(), "timeout must be >= 0")
	}
}
