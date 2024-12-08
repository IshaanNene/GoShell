#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Print commands and their arguments as they are executed
set -x

# Define the root directory of the project
ROOT_DIR=$(pwd)

# Define the test directory
TEST_DIR="$ROOT_DIR/tests"

# Create the test directory if it doesn't exist
mkdir -p $TEST_DIR

# Change to the test directory
cd $TEST_DIR

# Run the test command
go test $ROOT_DIR/...
