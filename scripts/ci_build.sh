#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Print commands and their arguments as they are executed
set -x

# Define the root directory of the project
ROOT_DIR=$(pwd)

# Define the build directory
BUILD_DIR="$ROOT_DIR/build"

# Create the build directory if it doesn't exist
mkdir -p $BUILD_DIR

# Change to the build directory
cd $BUILD_DIR

# Run the build command
go build -o myapp $ROOT_DIR/cmd/myapp

# Run tests
go test $ROOT_DIR/...
