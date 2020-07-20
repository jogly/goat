#!/usr/bin/env bash

DOCKER_BUILDKIT=1 \
  docker build \
  --build-arg BUILDKIT_INLINE_CACHE=1 \
  -t 727419036083.dkr.ecr.us-west-1.amazonaws.com/goat:cache \
  --target prod
  .

aws ecr get-login-password \
  --region us-west-1 \
  | docker login \
  --username AWS \
  --password-stdin \
  727419036083.dkr.ecr.us-west-1.amazonaws.com

docker push 727419036083.dkr.ecr.us-west-1.amazonaws.com/goat:cache
