# tagc
 tag your cli commands for future search

### Help
```bash
$ tagc --help

Usage of tagc:
-c string
command to store
-s    perform search with tags
-t string
tags to store separated by ','
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
