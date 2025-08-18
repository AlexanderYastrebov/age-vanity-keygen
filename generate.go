package main

import (
	"context"
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AlexanderYastrebov/age-vanity-keygen/internal/bech32"
	"github.com/AlexanderYastrebov/vanity25519"
)

type x25519identity struct {
	key      *ecdh.PrivateKey
	prefix   string
	elapsed  time.Duration
	attempts uint64
}

// Recipient returns the public X25519Recipient value corresponding to i.
func (i *x25519identity) Recipient() string {
	s, _ := bech32.Encode("age", i.key.PublicKey().Bytes())
	return s
}

// String returns the Bech32 private key encoding of i.
func (i *x25519identity) String() string {
	s, _ := bech32.Encode("AGE-SECRET-KEY-", i.key.Bytes())
	return strings.ToUpper(s)
}

// generateX25519Identity randomly generates a new X25519Identity with recipient prefix.
// prefix can not contain character '1'. It is transformed to lower case,
// characters 'b', 'i' and 'o' are replaced for '6', '7' and '0' respectively.
func generateX25519Identity(prefix string) (*x25519identity, error) {
	if len(prefix) == 0 {
		return nil, fmt.Errorf("empty prefix")
	}

	prefix = strings.ToLower(prefix)
	prefix = strings.TrimPrefix(prefix, "age1")
	if strings.Contains(prefix, "1") {
		return nil, fmt.Errorf("prefix can not contain \"1\"")
	}
	// map bech32 missing [bio] into [670]
	prefix = strings.ReplaceAll(prefix, "b", "6")
	prefix = strings.ReplaceAll(prefix, "i", "7")
	prefix = strings.ReplaceAll(prefix, "o", "0")
	prefix = "age1" + prefix

	hrp, prefixBytes, bits, err := bech32.DecodePrefix(prefix)
	if err != nil {
		return nil, fmt.Errorf("malformed public key prefix: %w", err)
	}
	if hrp != "age" {
		return nil, fmt.Errorf("malformed public key prefix human-readable part: %q", hrp)
	}

	k, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	offset, attempts := searchParallel(k.PublicKey().Bytes(), vanity25519.HasPrefixBits(prefixBytes, bits))
	elapsed := time.Since(start)

	vk, err := vanity25519.Add(k.Bytes(), offset)
	if err != nil {
		return nil, err
	}
	k, _ = ecdh.X25519().NewPrivateKey(vk)
	if err != nil {
		return nil, err
	}
	return &x25519identity{k, prefix, elapsed, attempts}, nil
}

func searchParallel(startPublicKey []byte, test func([]byte) bool) (*big.Int, uint64) {
	var result atomic.Pointer[big.Int]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var attemptsTotal atomic.Uint64
	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Go(func() {
			attempts := vanity25519.Search(ctx, startPublicKey, randBigInt(), 4096, test, func(_ []byte, offset *big.Int) {
				if result.CompareAndSwap(nil, offset) {
					cancel()
				}
			})
			attemptsTotal.Add(attempts)
		})
	}
	wg.Wait()

	return result.Load(), attemptsTotal.Load()
}

func randBigInt() *big.Int {
	var buf [8]byte
	rand.Read(buf[:])
	return new(big.Int).SetBytes(buf[:])
}
