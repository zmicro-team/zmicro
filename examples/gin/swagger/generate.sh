#!/bin/sh

echo 'Generating api'
PROTOS=$(find ./api -type f -name '*.proto')

for PROTO in $PROTOS; do
  echo $PROTO
  protoc \
    -I. -I$(dirname $PROTO) \
    -I../../third_party \
    --gofast_out=. \
    --gofast_opt paths=source_relative \
    --zmicro-gin_out=. \
    --zmicro-gin_opt paths=source_relative \
    --zmicro-gin_opt allow_empty_patch_body=true \
    $PROTO
done

echo 'Generating api swagger'
protoc \
  -I. \
  -I../../third_party \
  --openapiv2_out docs \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt allow_merge=true \
  --openapiv2_opt merge_file_name=swagger \
  --openapiv2_opt enums_as_ints=true \
  --openapiv2_opt json_names_for_fields=false \
$PROTOS