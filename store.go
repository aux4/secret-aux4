package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const storeVersion = 1

// Item is one stored secret. Values in Fields are ciphertext; nothing in this
// file is ever written in the clear.
type Item struct {
	Category string            `json:"category,omitempty"`
	Fields   map[string]string `json:"fields"`
}

// Store is the on-disk shape: vault -> item -> fields.
//
// The vault level is a namespace, not a physical container. Keeping it purely
// nominal is what lets a `secret://aux4/<vault>/<item>/<field>` reference in a
// shared config resolve on a machine that has never heard of that vault.
type Store struct {
	Version   int                         `json:"version"`
	Encrypter string                      `json:"encrypter,omitempty"`
	Vaults    map[string]map[string]*Item `json:"vaults"`
}

func newStore(encrypter string) *Store {
	return &Store{
		Version:   storeVersion,
		Encrypter: encrypter,
		Vaults:    map[string]map[string]*Item{},
	}
}

// splitRef splits a reference into vault and item. The vault is the first
// segment and the item is everything after it, matching how the aux4 core
// splits a secret:// URI (provider / ref... / field) and how the 1Password
// provider reads the same reference.
func splitRef(ref string) (string, string, error) {
	ref = strings.Trim(strings.TrimSpace(ref), "/")
	if ref == "" {
		return "", "", fmt.Errorf("ref is required, in the form <vault>/<item>")
	}

	parts := strings.SplitN(ref, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("ref %q must be in the form <vault>/<item>", ref)
	}

	return parts[0], parts[1], nil
}

func loadStore(path string, encrypter string) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return newStore(encrypter), nil
		}
		return nil, fmt.Errorf("could not read the store: %v", err)
	}

	store := newStore(encrypter)
	if err := json.Unmarshal(data, store); err != nil {
		return nil, fmt.Errorf("the store at %s is not valid JSON: %v", path, err)
	}
	if store.Vaults == nil {
		store.Vaults = map[string]map[string]*Item{}
	}

	return store, nil
}

// save writes the store atomically. A half-written store would lose every
// secret in it, so the rename only happens once the full contents are on disk.
func (s *Store) save(path string) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, data, 0600); err != nil {
		return fmt.Errorf("could not write the store: %v", err)
	}
	if err := os.Rename(temporary, path); err != nil {
		os.Remove(temporary)
		return fmt.Errorf("could not save the store: %v", err)
	}

	return nil
}

func (s *Store) get(vault string, item string) (*Item, bool) {
	items, ok := s.Vaults[vault]
	if !ok {
		return nil, false
	}
	entry, ok := items[item]
	return entry, ok
}

func (s *Store) put(vault string, item string, entry *Item) {
	if s.Vaults[vault] == nil {
		s.Vaults[vault] = map[string]*Item{}
	}
	s.Vaults[vault][item] = entry
}

func (s *Store) remove(vault string, item string) bool {
	items, ok := s.Vaults[vault]
	if !ok {
		return false
	}
	if _, ok := items[item]; !ok {
		return false
	}

	delete(items, item)
	if len(items) == 0 {
		delete(s.Vaults, vault)
	}
	return true
}

// reference is the secret:// URI for an entry, ready to paste into a config
// file or an .aux4.
func reference(vault string, item string, field string) string {
	if field == "" {
		return fmt.Sprintf("secret://aux4/%s/%s", vault, item)
	}
	return fmt.Sprintf("secret://aux4/%s/%s/%s", vault, item, field)
}

// entries returns every item in a stable order, optionally limited to one
// vault. Sorting matters: the output is a listing users read and diff.
func (s *Store) entries(vault string) []storeEntry {
	result := []storeEntry{}

	vaults := make([]string, 0, len(s.Vaults))
	for name := range s.Vaults {
		if vault != "" && name != vault {
			continue
		}
		vaults = append(vaults, name)
	}
	sort.Strings(vaults)

	for _, vaultName := range vaults {
		items := make([]string, 0, len(s.Vaults[vaultName]))
		for name := range s.Vaults[vaultName] {
			items = append(items, name)
		}
		sort.Strings(items)

		for _, itemName := range items {
			result = append(result, storeEntry{
				Vault: vaultName,
				Item:  itemName,
				Entry: s.Vaults[vaultName][itemName],
			})
		}
	}

	return result
}

type storeEntry struct {
	Vault string
	Item  string
	Entry *Item
}

// fieldNames returns an entry's field names in a stable order.
func (e storeEntry) fieldNames() []string {
	names := make([]string, 0, len(e.Entry.Fields))
	for name := range e.Entry.Fields {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
