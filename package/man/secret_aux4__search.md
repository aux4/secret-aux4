#### Description

The `search` command finds stored secrets whose item title contains the query, case-insensitively, and prints them as `secret://` references.

Only titles are searched. Field values are encrypted and are never scanned.

#### Usage

```bash
aux4 secret aux4 search <query> [--vault <name>] [--withFields <true|false>] [--store <path>]
```

query          The text to look for in item titles (required)
--vault        Limit the search to one vault (default: all vaults)
--withFields   Print one line per field instead of one per item (default: `false`)
--store        Path to the encrypted store (default: `${aux4HomeDir}/secret/aux4/store.json`)

#### Example

```bash
aux4 secret aux4 search billing
```

```text
secret://aux4/Work/Billing API
```
