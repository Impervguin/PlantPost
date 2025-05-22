#!/bin/bash

python3 scripts/minio-sync.py \
    --db-url "postgresql://impi:impi@localhost/plantpost" \
    --minio-endpoint "localhost:9000" \
    --minio-login "impi" \
    --minio-password "impichek" \
    --buckets "plants" "posts" \
    --fs-root ./deployments/data/ \
     --min-free-space 1024
