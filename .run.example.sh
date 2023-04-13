#!/bin/sh

cd "$(dirname "$0")"

go run ./... \
  --host=http://demo.testflo.com/server \
  --username=testflo \
  --password=testflo \
  --since='30 days ago' \
  "$@"
