# LxRoot New Website

### How do i add Collection?
Add collection manually in Community Edition 7.6.1 build 3200 (as REST API request only allowed in enterprise edition)

## Bucket Export
> cbexport json -c couchbase://127.0.0.1 -u Administrator -p Mostain321$  -b lxroot -o data.json -f lines -t 4 --scope-field scope --collection-field collection

## Bucket Import
> cbimport json -c couchbase://127.0.0.1 -u Administrator -p Mostain321$ -b lxerp -d file://data.json -f lines -g %id% --scope-collection-exp %scope%.%collection%

### Reference
* https://docs.couchbase.com/server/current/tools/cbexport-json.html
* https://docs.couchbase.com/server/current/tools/cbimport-json.html