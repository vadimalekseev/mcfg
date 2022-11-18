install:

`go install github.com/vadimalekseev/mcfg@v0.0.1`

usage:

`mcfg merge -p ".yaml$" -f staging > .o3/k8s/values_staging.yaml`
`mcfg merge -p ".yaml$" -f production > .o3/k8s/values_production.yaml`
