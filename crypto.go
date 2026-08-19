package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// The crypto is not implemented here. aux4/encrypter already provides it behind
// a provider facade, so this package only decides *what* to encrypt and where to
// keep it -- and inherits aws-kms for free by passing --provider through.
const defaultEncrypter = "aux4"

// runEncrypter invokes the encrypter facade with the plaintext (or ciphertext)
// on stdin and the key in the environment.
//
// Both matter: an argument list is world-readable through `ps`, so neither the
// secret being encrypted nor the key that protects it is ever passed as one.
// AUX4_SECRET is the env binding aux4/encrypter already declares for --secret.
func runEncrypter(action string, encrypter string, key string, input string) (string, error) {
	if encrypter == "" {
		encrypter = defaultEncrypter
	}

	command := exec.Command("aux4", "encrypter", action, "--provider", encrypter)
	command.Stdin = strings.NewReader(input)
	command.Env = append(os.Environ(), "AUX4_SECRET="+key)

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return "", fmt.Errorf("encrypter %s failed: %s", action, message)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func encryptValue(encrypter string, key string, plaintext string) (string, error) {
	return runEncrypter("encrypt", encrypter, key, plaintext)
}

func decryptValue(encrypter string, key string, ciphertext string) (string, error) {
	return runEncrypter("decrypt", encrypter, key, ciphertext)
}

// generateKey asks the encrypter for fresh key material.
func generateKey(encrypter string) (string, error) {
	if encrypter == "" {
		encrypter = defaultEncrypter
	}

	command := exec.Command("aux4", "encrypter", "generate-secret", "--provider", encrypter)

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return "", fmt.Errorf("could not generate a key: %s", message)
	}

	key := strings.TrimSpace(stdout.String())
	if key == "" {
		return "", fmt.Errorf("the encrypter returned an empty key")
	}
	return key, nil
}

// resolveKey finds the key that protects the store.
//
// AUX4_SECRET_KEY wins, which is what makes the store usable in CI and in
// containers where no interactive setup has run. Otherwise a key file sits next
// to the store, generated on first write.
//
// Worth being precise about what that buys: a key file beside the store is not
// meaningfully stronger than a 0600 plaintext file against anything running as
// this user. It protects the store in backups, in synced home directories, and
// from casual inspection. Set AUX4_SECRET_KEY from somewhere else for more.
func resolveKey(storePath string, encrypter string, create bool) (string, error) {
	if key := os.Getenv("AUX4_SECRET_KEY"); key != "" {
		return key, nil
	}

	keyPath := filepath.Join(filepath.Dir(storePath), "key")

	data, err := os.ReadFile(keyPath)
	if err == nil {
		key := strings.TrimSpace(string(data))
		if key != "" {
			return key, nil
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("could not read the key file: %v", err)
	}

	if !create {
		return "", fmt.Errorf("no key found at %s and AUX4_SECRET_KEY is not set", keyPath)
	}

	key, err := generateKey(encrypter)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(keyPath), 0700); err != nil {
		return "", err
	}
	if err := os.WriteFile(keyPath, []byte(key), 0600); err != nil {
		return "", fmt.Errorf("could not write the key file: %v", err)
	}

	return key, nil
}
