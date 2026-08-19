# secret aux4 remove

## deleting a secret

```beforeAll
aux4 secret aux4 create --vault Work --item Jira --fields "token=a" --store .secret/store.json > /dev/null 2>&1
aux4 secret aux4 create --vault Work --item Confluence --fields "token=b" --store .secret/store.json > /dev/null 2>&1
```

```afterAll
rm -rf .secret
```

### should delete the item and leave the others

```execute
aux4 secret aux4 remove --ref Work/Jira --store .secret/store.json
aux4 secret aux4 list --store .secret/store.json
```

```expect
secret://aux4/Work/Jira removed
secret://aux4/Work/Confluence
```

### should drop the vault once its last item is gone

A vault is only a namespace, so it stops existing when nothing is left in it.

```execute
aux4 secret aux4 remove --ref Work/Jira --store .secret/store.json > /dev/null
aux4 secret aux4 remove --ref Work/Confluence --store .secret/store.json > /dev/null
aux4 secret aux4 list --store .secret/store.json
echo "(store empty)"
```

```expect
(store empty)
```

### should fail on an unknown item

```execute
aux4 secret aux4 remove --ref Work/Nope --store .secret/store.json
```

```error:partial
no secret found at secret://aux4/Work/Nope
```
