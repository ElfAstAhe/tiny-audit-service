#!/usr/bin/env bash

cd ..

mkdir -p api/{proto/tiny-audit-service,rest} \
    cmd/tiny-audit-service \
    configs \
    deployments \
    docs \
    internal/{app,config,domain,facade/{dto,mapper},repository/postgres,transport/{errs,rest,grpc},usecase} \
    migrations/tiny-audit-service \
    pkg \
    scripts

cd scripts
