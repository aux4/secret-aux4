#### Description

The `list` command prints stored secrets as `secret://` references, one per line, so they can be copied straight into a config file or an `.aux4`.

Listing does not decrypt anything. Vault names, item titles and field names are stored in the clear precisely so the store can be browsed without unlocking it — which is also why none of them should hold anything sensitive.

Output is sorted by vault then item, so it is stable between runs and diffable.

#### Usage

```bash
aux4 secret aux4 list [--vault <name>] [--withFields <true|false>] [--store <path>]
```

--vault        Limit the listing to one vault (default: all vaults)
--withFields   Print one line per field instead of one per item (default: `false`)
--store        Path to the encrypted store (default: `${aux4HomeDir}/secret/aux4/store.json`)

#### Example

```bash
aux4 secret aux4 list
```

```text
secret://aux4/Personal/GitHub
secret://aux4/Work/Billing API
```

```bash
aux4 secret aux4 list --vault Work --withFields true
```

```text
secret://aux4/Work/Billing API/clientId
secret://aux4/Work/Billing API/clientSecret
```
