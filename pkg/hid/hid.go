package hid

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/Palladium-blockchain/go-human-readable-id/internal/words"
)

var (
	ErrUnknownGenerator = errors.New("unknown generator")
	ErrUnclosedToken    = errors.New("unclosed token")
)

func Generate(template string, opts ...Option) (string, error) {
	return GenerateContext(context.Background(), template, opts...)
}

func GenerateContext(ctx context.Context, template string, opts ...Option) (string, error) {
	cfg := Config{
		Strict:     true,
		Generators: nil,
		Rand: rand.New(
			rand.NewPCG(
				uint64(time.Now().UnixNano()),
				0x9e3779b97f4a7c15,
			),
		),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.Generators == nil {
		cfg.Generators = make(map[string]Generator)
	}

	s := template
	var resBuf strings.Builder
	const (
		stateSym = iota
		stateKey
	)
	state := stateSym
	var keyBuf strings.Builder
	tokenStart := -1
	for i := 0; i < len(s); {
		switch state {
		case stateSym:
			if s[i] != '{' {
				resBuf.WriteByte(s[i])
				i++
				continue
			}
			if s[i] == '{' {
				tokenStart = i
				state = stateKey
				i++
				continue
			}
		case stateKey:
			if s[i] != '}' {
				keyBuf.WriteByte(s[i])
				i++
				continue
			}
			if s[i] == '}' {
				state = stateSym
				i++

				key := keyBuf.String()
				keyBuf.Reset()
				gen, ok := cfg.Generators[key]
				if !ok {
					if cfg.Strict {
						return "", fmt.Errorf("unknown generator %s as %d pos: %w", key, tokenStart, ErrUnknownGenerator)
					}
					resBuf.WriteByte('{')
					resBuf.WriteString(key)
					resBuf.WriteByte('}')
					continue
				}

				out, err := gen(ctx, &cfg)
				if err != nil {
					return "", fmt.Errorf("generator %s failed at %d: %w", key, tokenStart, err)
				}

				resBuf.WriteString(out)
				continue
			}
		}
	}
	if state == stateKey {
		if cfg.Strict {
			return "", fmt.Errorf("unclosed token at %d: %w", tokenStart, ErrUnclosedToken)
		}
		resBuf.WriteByte('{')
		resBuf.WriteString(keyBuf.String())
	}

	return resBuf.String(), nil
}

type Config struct {
	Generators map[string]Generator
	Strict     bool
	Rand       *rand.Rand
}

type Generator func(ctx context.Context, cfg *Config) (string, error)

type Option func(*Config)

func WithDefaultGenerators() Option {
	return func(c *Config) {
		if c.Generators == nil {
			c.Generators = make(map[string]Generator)
		}
		for k, v := range DefaultGenerators() {
			if _, exists := c.Generators[k]; !exists {
				c.Generators[k] = v
			}
		}
	}
}

func WithGenerator(token string, generator Generator) Option {
	return func(c *Config) {
		if c.Generators == nil {
			c.Generators = make(map[string]Generator)
		}
		c.Generators[token] = generator
	}
}

func WithSeed(seed uint64) Option {
	return func(c *Config) {
		c.Rand = rand.New(rand.NewPCG(seed, 0x9e3779b97f4a7c15))
	}
}

func WithStrict(strict bool) Option {
	return func(c *Config) { c.Strict = strict }
}

func AdjGenerator(ctx context.Context, cfg *Config) (string, error) {
	return WordGenerator(words.AdjWords)(ctx, cfg)
}

func NounGenerator(ctx context.Context, cfg *Config) (string, error) {
	return WordGenerator(words.NounWords)(ctx, cfg)
}

func VerbGenerator(ctx context.Context, cfg *Config) (string, error) {
	return WordGenerator(words.VerbWords)(ctx, cfg)
}

func WordGenerator(words []string) Generator {
	return func(ctx context.Context, cfg *Config) (string, error) {
		if len(words) == 0 {
			return "", errors.New("empty word list")
		}
		i := cfg.Rand.IntN(len(words))
		return words[i], nil
	}
}

func IntGenerator(min, max int) Generator {
	return func(ctx context.Context, cfg *Config) (string, error) {
		if min > max {
			return "", errors.New("min > max")
		}
		n := cfg.Rand.IntN(max-min+1) + min
		return strconv.Itoa(n), nil
	}
}

func DefaultGenerators() map[string]Generator {
	return map[string]Generator{
		"adj":     AdjGenerator,
		"noun":    NounGenerator,
		"verb":    VerbGenerator,
		"digit":   IntGenerator(0, 9),
		"2-digit": IntGenerator(10, 99),
		"3-digit": IntGenerator(100, 999),
	}
}
