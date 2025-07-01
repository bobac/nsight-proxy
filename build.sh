#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Define target platforms (OS and Architecture combinations)
PLATFORMS=("windows/amd64" "windows/arm64" "linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64")

# Define commands to build (corresponds to directories in cmd/)
COMMANDS=("getdata" "fetchall" "nsight-proxy")

# Output directory
OUTPUT_DIR="bin"

# Get the package path (assuming the script is run from the project root)
# This might need adjustment depending on your go module name if it differs significantly
PACKAGE_PATH=$(go list -m)
echo "Using package path: $PACKAGE_PATH"

echo "Starting build process..."

# Loop through each command
for cmd in "${COMMANDS[@]}"; do
    echo "Building command: $cmd ..."
    # Loop through each platform
    for platform in "${PLATFORMS[@]}"; do
        # Split platform into OS and Arch
        IFS='/' read -r GOOS GOARCH <<< "$platform"

        # Set output path
        output_name="$OUTPUT_DIR/${GOOS}_${GOARCH}/${cmd}"

        # Add .exe suffix for windows
        if [ "$GOOS" = "windows" ]; then
            output_name+=".exe"
        fi

        # Ensure the target directory exists
        mkdir -p "$(dirname "$output_name")"

        # Build the command with ldflags directly specified
        echo "  Building for $GOOS/$GOARCH -> $output_name"
        env GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-s -w" -o "$output_name" "${PACKAGE_PATH}/cmd/${cmd}"
        if [ $? -ne 0 ]; then
            echo "  Build FAILED for $cmd on $GOOS/$GOARCH"
        else
            echo "  Build SUCCESSFUL for $cmd on $GOOS/$GOARCH"
        fi
    done
done

echo "Build process completed."

# Make the script executable (optional, user might need to do this)
# chmod +x build.sh
echo "Run 'chmod +x build.sh' to make the script executable." 