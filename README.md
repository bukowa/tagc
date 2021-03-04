# tagc
tag any text for future search

### Help
```bash
$ tagc --help

Usage of tagc:
-c string text to store
-s bool perform search with tags
-t string tags to store separated by ','
```

### Examples
````bash
$ tagc -t 1,2,3 -c "my command"
$ tagc -t 2,3 -c "other command"
$ tagc -s -t 2,3
{
 "Matches": {
  "2": [
   "my command",
   "other command"
  ]
 }
}

````

#### todo
- include git based backup with new flag
