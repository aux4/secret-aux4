package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

func argAt(args []string, index int) string {
	if len(args) > index {
		return strings.TrimSpace(args[index])
	}
	return ""
}

// splitList reads a comma-separated flag value.
func splitList(raw string) []string {
	result := []string{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// parseFields reads "key=value,key=value" pairs. Splitting each pair on the
// first = only means a value may itself contain =, which secrets frequently do
// (base64 padding, connection strings).
func parseFields(raw string) (map[string]string, []string, error) {
	fields := map[string]string{}
	order := []string{}

	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		index := strings.Index(pair, "=")
		if index <= 0 {
			return nil, nil, fmt.Errorf("field %q must be in the form name=value", pair)
		}

		name := strings.TrimSpace(pair[:index])
		if name == "" {
			return nil, nil, fmt.Errorf("field %q has an empty name", pair)
		}

		fields[name] = pair[index+1:]
		order = append(order, name)
	}

	if len(fields) == 0 {
		return nil, nil, fmt.Errorf("at least one field is required, in the form name=value")
	}

	return fields, order, nil
}

// args: storePath encrypter ref fields otp
//
// This is the command the aux4 core calls to resolve a secret:// reference, so
// its contract is fixed: a JSON object of field name to plaintext on stdout.
func runGet(args []string) {
	storePath := argAt(args, 0)
	encrypter := argAt(args, 1)
	ref := argAt(args, 2)
	fields := splitList(argAt(args, 3))

	vault, item, err := splitRef(ref)
	if err != nil {
		fail("%v", err)
	}
	if len(fields) == 0 {
		fail("fields is required")
	}

	store, err := loadStore(storePath, encrypter)
	if err != nil {
		fail("%v", err)
	}

	entry, ok := store.get(vault, item)
	if !ok {
		fail("no secret found at %s", reference(vault, item, ""))
	}

	key, err := resolveKey(storePath, encrypter, false)
	if err != nil {
		fail("%v", err)
	}

	// Ordered so the output is stable and diffable; the aux4 core reads it as a
	// map, but a human reading `secret aux4 get` should get the same bytes twice.
	result := map[string]string{}
	names := []string{}

	for _, field := range fields {
		ciphertext, ok := entry.Fields[field]
		if !ok {
			fail("field %q not found in %s", field, reference(vault, item, ""))
		}

		plaintext, err := decryptValue(store.Encrypter, key, ciphertext)
		if err != nil {
			fail("%v", err)
		}

		result[field] = plaintext
		names = append(names, field)
	}

	sort.Strings(names)
	ordered := make([]string, 0, len(names))
	for _, name := range names {
		encoded, _ := json.Marshal(result[name])
		nameEncoded, _ := json.Marshal(name)
		ordered = append(ordered, fmt.Sprintf("  %s: %s", nameEncoded, encoded))
	}

	fmt.Fprintf(os.Stdout, "{\n%s\n}\n", strings.Join(ordered, ",\n"))
}

// args: storePath encrypter vault item fields category
func runCreate(args []string) {
	storePath := argAt(args, 0)
	encrypter := argAt(args, 1)
	vault := argAt(args, 2)
	item := argAt(args, 3)
	rawFields := argAt(args, 4)
	category := argAt(args, 5)

	if vault == "" {
		fail("vault is required")
	}
	if item == "" {
		fail("item is required")
	}
	if strings.Contains(vault, "/") {
		fail("vault %q must not contain '/', which separates the parts of a reference", vault)
	}

	fields, _, err := parseFields(rawFields)
	if err != nil {
		fail("%v", err)
	}

	store, err := loadStore(storePath, encrypter)
	if err != nil {
		fail("%v", err)
	}
	if _, exists := store.get(vault, item); exists {
		fail("%s already exists, use set to change a field", reference(vault, item, ""))
	}

	key, err := resolveKey(storePath, encrypter, true)
	if err != nil {
		fail("%v", err)
	}

	entry := &Item{Category: category, Fields: map[string]string{}}
	for name, value := range fields {
		ciphertext, err := encryptValue(store.Encrypter, key, value)
		if err != nil {
			fail("%v", err)
		}
		entry.Fields[name] = ciphertext
	}

	store.put(vault, item, entry)
	if err := store.save(storePath); err != nil {
		fail("%v", err)
	}

	fmt.Fprintln(os.Stdout, reference(vault, item, ""))
}

// args: storePath encrypter ref field value
func runSet(args []string) {
	storePath := argAt(args, 0)
	encrypter := argAt(args, 1)
	ref := argAt(args, 2)
	field := argAt(args, 3)
	value := ""
	if len(args) > 4 {
		value = args[4]
	}

	vault, item, err := splitRef(ref)
	if err != nil {
		fail("%v", err)
	}
	if field == "" {
		fail("field is required")
	}

	store, err := loadStore(storePath, encrypter)
	if err != nil {
		fail("%v", err)
	}

	entry, ok := store.get(vault, item)
	if !ok {
		// Creating on set would silently produce a half-populated item under a
		// mistyped ref, which then fails much later at resolution time.
		fail("no secret found at %s, use create first", reference(vault, item, ""))
	}

	key, err := resolveKey(storePath, encrypter, true)
	if err != nil {
		fail("%v", err)
	}

	ciphertext, err := encryptValue(store.Encrypter, key, value)
	if err != nil {
		fail("%v", err)
	}

	entry.Fields[field] = ciphertext
	if err := store.save(storePath); err != nil {
		fail("%v", err)
	}

	fmt.Fprintf(os.Stdout, "%s updated\n", reference(vault, item, field))
}

// args: storePath encrypter vault withFields
func runList(args []string) {
	storePath := argAt(args, 0)
	encrypter := argAt(args, 1)
	vault := argAt(args, 2)
	withFields := argAt(args, 3) == "true"

	store, err := loadStore(storePath, encrypter)
	if err != nil {
		fail("%v", err)
	}

	printEntries(store.entries(vault), withFields)
}

// args: storePath encrypter query vault withFields
func runSearch(args []string) {
	storePath := argAt(args, 0)
	encrypter := argAt(args, 1)
	query := argAt(args, 2)
	vault := argAt(args, 3)
	withFields := argAt(args, 4) == "true"

	if query == "" {
		fail("query is required")
	}

	store, err := loadStore(storePath, encrypter)
	if err != nil {
		fail("%v", err)
	}

	needle := strings.ToLower(query)
	matched := []storeEntry{}
	for _, entry := range store.entries(vault) {
		if strings.Contains(strings.ToLower(entry.Item), needle) {
			matched = append(matched, entry)
		}
	}

	printEntries(matched, withFields)
}

// printEntries writes secret:// references, one per line, so the output can be
// pasted straight into a config file.
func printEntries(entries []storeEntry, withFields bool) {
	for _, entry := range entries {
		if !withFields {
			fmt.Fprintln(os.Stdout, reference(entry.Vault, entry.Item, ""))
			continue
		}
		for _, field := range entry.fieldNames() {
			fmt.Fprintln(os.Stdout, reference(entry.Vault, entry.Item, field))
		}
	}
}

// args: storePath encrypter ref
func runRemove(args []string) {
	storePath := argAt(args, 0)
	encrypter := argAt(args, 1)
	ref := argAt(args, 2)

	vault, item, err := splitRef(ref)
	if err != nil {
		fail("%v", err)
	}

	store, err := loadStore(storePath, encrypter)
	if err != nil {
		fail("%v", err)
	}

	if !store.remove(vault, item) {
		fail("no secret found at %s", reference(vault, item, ""))
	}
	if err := store.save(storePath); err != nil {
		fail("%v", err)
	}

	fmt.Fprintf(os.Stdout, "%s removed\n", reference(vault, item, ""))
}
