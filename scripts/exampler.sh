#!/bin/bash

python3 ./scripts/exampler.py --dirs ./config ./deployments \
    --extensions .yaml .env \
    --suffix .example \
    --exclude docker-compose traefik.yaml dynamic.yaml \
    --override