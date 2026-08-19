# secret aux4 search

## finding secrets by title

```beforeAll
aux4 secret aux4 create --vault Personal --item GitHub --fields "password=a" --store .secret/store.json > /dev/null 2>&1
aux4 secret aux4 create --vault Work --item "Billing API" --fields "clientSecret=b" --store .secret/store.json > /dev/null 2>&1
```

```afterAll
rm -rf .secret
```

### should match part of a title

```execute
aux4 secret aux4 search billing --store .secret/store.json
```

```expect
secret://aux4/Work/Billing API
```

### should ignore case

```execute
aux4 secret aux4 search GITHUB --store .secret/store.json
```

```expect
secret://aux4/Personal/GitHub
```

### should print nothing when nothing matches

```execute
aux4 secret aux4 search nothinghere --store .secret/store.json
echo "(no matches)"
```

```expect
(no matches)
```
