#!/bin/sh

echo 'Generating api'
PROTOS=$(find ./api -type f -name '*.proto')

for PROTO in $PROTOS; do
  echo $PROTO
  protoc \
    -I. -I$(dirname $PROTO) \
    -I../third_party \
    --gofast_out=. \
    --gofast_opt paths=source_relative \
    --zmicro-gin_out=. \
    --zmicro-gin_opt paths=source_relative \
    --zmicro-gin_opt allow_empty_patch_body=true \
    $PROTO
done

echo 'Generating errno'
ERRORS=$(find ./errno -type f -name '*.proto')
for ERROR in $ERRORS; do
  echo $ERROR
  protoc \
  -I. -I../..\
  --gofast_out=. \
  --gofast_opt paths=source_relative \
  --zmicro-errno_out=. \
  --zmicro-errno_opt epk=github.com/zmicro-team/zmicro/core/errors \
  --zmicro-errno_opt paths=source_relative \
  $ERROR
done
