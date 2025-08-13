# CLI Testing Justfile

# Default recipe - run all tests
default: test-basic test-compose

# Test basic CLI functionality
test-basic:
    go run . -version

# Test with basic compose file
test-compose:
    go run . -f tests/basic-test.yaml

# Test with complex compose file
test-complex:
    go run . -f tests/complex-test.yaml

# Test variable interpolation using .env file in tests directory
test-env:
    go run . -f tests/interpolation-test.yaml -env-file tests/.env -dry-run

# Clean output directory
clean:
    rm -rf tests/output/*

# Run all tests and clean up
test-all: test-basic test-compose test-complex test-env clean