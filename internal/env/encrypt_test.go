package env

import (
	"strings"
	"testing"
)

func TestEncryptor_RoundTrip(t *testing.T) {
	e := NewEncryptor("supersecret")
	plaintext := "my-vault-secret-value"

	encoded, err := e.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: unexpected error: %v", err)
	}
	if encoded == plaintext {
		t.Fatal("Encrypt: output should not equal plaintext")
	}

	decoded, err := e.Decrypt(encoded)
	if err != nil {
		t.Fatalf("Decrypt: unexpected error: %v", err)
	}
	if decoded != plaintext {
		t.Fatalf("Decrypt: got %q, want %q", decoded, plaintext)
	}
}

func TestEncryptor_DifferentPassphrasesFail(t *testing.T) {
	e1 := NewEncryptor("passphrase-one")
	e2 := NewEncryptor("passphrase-two")

	encoded, err := e1.Encrypt("secret")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	_, err = e2.Decrypt(encoded)
	if err == nil {
		t.Fatal("Decrypt: expected error with wrong passphrase")
	}
}

func TestEncryptor_NonceDiffersEachCall(t *testing.T) {
	e := NewEncryptor("stable-key")
	a, _ := e.Encrypt("value")
	b, _ := e.Encrypt("value")
	if a == b {
		t.Fatal("Encrypt: expected different ciphertext each call due to random nonce")
	}
}

func TestEncryptor_Decrypt_InvalidBase64(t *testing.T) {
	e := NewEncryptor("key")
	_, err := e.Decrypt("!!!not-base64!!!")
	if err == nil {
		t.Fatal("Decrypt: expected error for invalid base64")
	}
}

func TestEncryptor_Decrypt_TooShort(t *testing.T) {
	e := NewEncryptor("key")
	import64 := "YQ=="olean // base64 of "a" — too short for nonce
	_ = import64
	// encode a single byte
	import "encoding/base64"
	tiny := base64.StdEncoding.EncodeToString([]byte(strings.Repeat("x", 3)))
	_, err := e.Decrypt(tiny)
	if err == nil {
		t.Fatal("Decrypt: expected error for too-short ciphertext")
	}
}
