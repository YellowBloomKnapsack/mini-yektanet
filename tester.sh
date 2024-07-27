#!/bin/bash

# Function to check if a directory contains test files
contains_test_files() {
    if ls "$1"/*_test.go &> /dev/null; then
        return 0
    else
        return 1
    fi
}

RED='\033[0;31m'
PURPLE='\033[1;35m'
GREEN='\033[0;32m'
BLUE='\033[1;34m'
NC='\033[0m'

# Function to run tests and log coverage with colored output
run_tests_with_coverage() {
    for dir in $(find . -type d); do
        if contains_test_files "$dir"; then
            echo -e "${BLUE}- Running tests in $dir ${NC}"

            # Run the tests and capture the output
            result=$(go test -coverprofile=coverage.out -v $dir 2>&1)

            # Print each test function result with color
            echo "$result" | while IFS= read -r line; do
                if [[ "$line" == *"--- PASS"* ]]; then
                    echo -e "${GREEN}✓ $line ${NC}"  # Green for pass
                elif [[ "$line" == *"--- FAIL"* ]]; then
                    echo -e "${RED}✕ $line ${NC}"  # Red for fail
                fi
            done

            # Print the coverage
            if [ -f coverage.out ]; then
                coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
                echo -e "${PURPLE} Coverage for ${BLUE}$dir${PURPLE}: $coverage\n ${NC}"
                echo -e "---------------------------------------------\n"
                rm coverage.out
            fi
        fi
    done
}

# Execute the function
run_tests_with_coverage
