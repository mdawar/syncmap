_default:
  @just -l -u

# Run the tests.
test *args:
  go test -race -cover -vet=all -count 1 {{args}}

