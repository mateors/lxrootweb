# LxRoot New Website

> go clean -cache

> go clean -modcache


```
SELECT OBJECT_NAMES(b) FROM lxroot._default.company b;

SELECT OBJECT_NAMES(b)as fields FROM lxroot._default.company b;

SELECT DISTINCT OBJECT_NAMES(b)as fields FROM lxroot._default.company b;

CREATE INDEX `idx_obj_name_db` ON `db`((object_names(db)))  USING GSI;
```


### System info query
* `SELECT * FROM system:datastores`
* `SELECT * FROM system:namespaces`
* `SELECT * FROM system:buckets`
* `SELECT * FROM system:scopes`
* `SELECT * FROM system:keyspaces` | list collection
* `SELECT * FROM system:indexes`


### How do i add Collection?
Add collection manually in Community Edition 7.6.1 build 3200 (as REST API request only allowed in enterprise edition)

## Bucket Export
> cbexport json -c couchbase://127.0.0.1 -u Administrator -p Mostain321$  -b lxroot -o data.json -f lines -t 4 --scope-field scope --collection-field collection

## Bucket Import
> cbimport json -c couchbase://127.0.0.1 -u Administrator -p Mostain321$ -b lxerp -d file://data.json -f lines -g %id% --scope-collection-exp %scope%.%collection%

### Reference
* https://docs.couchbase.com/server/current/tools/cbexport-json.html
* https://docs.couchbase.com/server/current/tools/cbimport-json.html
* https://docs.couchbase.com/server/current/n1ql/n1ql-intro/sysinfo.html