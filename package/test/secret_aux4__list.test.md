# secret aux4 list

## listing stored secrets

```beforeAll
aux4 secret aux4 create --vault Personal --item GitHub --fields "password=a" --store .secret/store.json > /dev/null 2>&1
aux4 secret aux4 create --vault Work --item "Billing API" --fields "clientId=b,clientSecret=c" --store .secret/store.json > /dev/null 2>&1
```

```afterAll
rm -rf .secret
```

### should print every reference sorted by vault and item

```execute
aux4 secret aux4 list --store .secret/store.json
```

```expect
secret://aux4/Personal/GitHub
secret://aux4/Work/Billing API
```

### should limit the listing to one vault

```execute
aux4 secret aux4 list --vault Work --store .secret/store.json
```

```expect
secret://aux4/Work/Billing API
```

### should print one line per field

```execute
aux4 secret aux4 list --vault Work --withFields true --store .secret/store.json
```

```expect
secret://aux4/Work/Billing API/clientId
secret://aux4/Work/Billing API/clientSecret
```

### should print nothing for an empty store

```execute
aux4 secret aux4 list --store .secret/empty.json
echo "(no output)"
```

```expect
(no output)
```
