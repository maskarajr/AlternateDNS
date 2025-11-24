#!/bin/bash
echo "Building AlternateDNS for Linux/macOS (portable)..."

CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -o AlternateDNS

if [ $? -eq 0 ]; then
    echo "Build successful! AlternateDNS created."
    echo "File size:"
    ls -lh AlternateDNS
else
    echo "Build failed!"
    exit 1
fi

