# secret aux4 set

## updating a field

```beforeAll
aux4 secret aux4 create --vault Work --item Jira --fields "token=original,user=john" --store .secret/store.json > /dev/null 2>&1
```

```afterAll
rm -rf .secret
```

### should replace the value

```execute
aux4 secret aux4 set --ref Work/Jira --field token --value rotated --store .secret/store.json
aux4 secret aux4 get --ref Work/Jira --fields token --store .secret/store.json
```

```expect
secret://aux4/Work/Jira/token updated
{
  "token": "rotated"
}
```

### should leave the other fields untouched

```execute
aux4 secret aux4 set --ref Work/Jira --field token --value rotated --store .secret/store.json > /dev/null
aux4 secret aux4 get --ref Work/Jira --fields user --store .secret/store.json
```

```expect:json
{
  "user": "john"
}
```

### should add a field that does not exist yet

```execute
aux4 secret aux4 set --ref Work/Jira --field extra --value added --store .secret/store.json > /dev/null
aux4 secret aux4 get --ref Work/Jira --fields extra --store .secret/store.json
```

```expect:json
{
  "extra": "added"
}
```

### should refuse to create the item

```execute
aux4 secret aux4 set --ref Work/Missing --field token --value v --store .secret/store.json
```

```error:partial
no secret found at secret://aux4/Work/Missing, use create first
```
