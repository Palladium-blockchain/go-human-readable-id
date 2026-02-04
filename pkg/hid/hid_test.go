package hid

import (
	"context"
	"errors"
	"testing"
)

func TestGenerateContext_ReplacesMultipleTokens(t *testing.T) {
	ctx := context.Background()

	got, err := GenerateContext(
		ctx,
		"{adj}%{num}(something).txt",
		WithGenerator("adj", func(context.Context, *Config) (string, error) { return "fluffy", nil }),
		WithGenerator("num", func(context.Context, *Config) (string, error) { return "42", nil }),
		WithStrict(true),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "fluffy%42(something).txt"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGenerateContext_UnknownKey_StrictErrors(t *testing.T) {
	ctx := context.Background()

	_, err := GenerateContext(
		ctx,
		"{adj}-{unknown}-{num}",
		WithGenerator("adj", func(context.Context, *Config) (string, error) { return "fluffy", nil }),
		WithGenerator("num", func(context.Context, *Config) (string, error) { return "42", nil }),
		WithStrict(true),
	)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrUnknownGenerator) {
		t.Fatalf("expected ErrUnknownGenerator, got: %v", err)
	}
}

func TestGenerateContext_UnknownKey_NonStrictKeepsToken(t *testing.T) {
	ctx := context.Background()

	got, err := GenerateContext(
		ctx,
		"{adj}-{unknown}-{num}",
		WithGenerator("adj", func(context.Context, *Config) (string, error) { return "fluffy", nil }),
		WithGenerator("num", func(context.Context, *Config) (string, error) { return "42", nil }),
		WithStrict(false),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "fluffy-{unknown}-42"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGenerateContext_UnclosedToken_StrictErrors(t *testing.T) {
	ctx := context.Background()

	_, err := GenerateContext(
		ctx,
		"prefix-{adj",
		WithGenerator("adj", func(context.Context, *Config) (string, error) { return "fluffy", nil }),
		WithStrict(true),
	)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrUnclosedToken) {
		t.Fatalf("expected ErrUnclosedToken, got: %v", err)
	}
}

func TestGenerateContext_UnclosedToken_NonStrictKeepsRemainder(t *testing.T) {
	ctx := context.Background()

	got, err := GenerateContext(
		ctx,
		"prefix-{adj",
		WithGenerator("adj", func(context.Context, *Config) (string, error) { return "fluffy", nil }),
		WithStrict(false),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// В non-strict режиме незакрытый токен остаётся как есть.
	want := "prefix-{adj"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWithDefaultGenerators_DoesNotOverrideExisting(t *testing.T) {
	ctx := context.Background()

	got, err := GenerateContext(
		ctx,
		"{adj}",
		WithGenerator("adj", func(context.Context, *Config) (string, error) { return "CUSTOM", nil }),
		WithDefaultGenerators(), // должен мерджить, а не перезатирать
		WithStrict(true),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "CUSTOM" {
		t.Fatalf("got %q, want %q", got, "CUSTOM")
	}
}
