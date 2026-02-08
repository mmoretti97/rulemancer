#!/usr/bin/env bash

# Spawn the binary, get the PID and store it in RULEMANCER_PID, and parse the output of the spawn to get the ADMIN token
# The admin token is written in stdout as: admin jwt: TOKEN
# The token is stored in API_TOKEN

# Check if the env RULEMANCER_JWT_SECRET exists
if [ -z "$RULEMANCER_JWT_SECRET" ]; then
    echo "Error: RULEMANCER_JWT_SECRET environment variable not set"
    return 1 2>/dev/null || exit 1
fi

# Check if rulemancer binary exists
BINARY="./rulemancer"
if [ ! -f "$BINARY" ]; then
    echo "Error: rulemancer binary not found at $BINARY"
    return 1 2>/dev/null || exit 1
fi

# Create a temporary file to capture output
TEMP_OUTPUT=$(mktemp)

# Start rulemancer in background and capture output
"$BINARY" serve > "$TEMP_OUTPUT" 2>&1 &

cat $TEMP_OUTPUT

# Get the PID
export RULEMANCER_PID=$!

echo "Rulemancer started with PID: $RULEMANCER_PID"

# Wait for the token to appear in output (max 10 seconds)
for i in {1..20}; do
    if grep -q "admin jwt:" "$TEMP_OUTPUT"; then
        break
    fi
    sleep 0.5
done

# Extract the admin token
export API_TOKEN=$(grep "admin jwt:" "$TEMP_OUTPUT" | sed 's/admin jwt: //' | tr -d '[:space:]')

if [ -z "$API_TOKEN" ]; then
    echo "Warning: Could not extract admin token"
    rm -f "$TEMP_OUTPUT"
    return 1 2>/dev/null || exit 1
fi

echo "Admin token extracted: ${API_TOKEN:0:20}..."
echo "Environment variables set:"
echo "  RULEMANCER_PID=$RULEMANCER_PID"
echo "  API_TOKEN=$API_TOKEN"

# Clean up temp file
rm -f "$TEMP_OUTPUT"

# Keep the binary running in background
echo "Rulemancer is running in background. Use 'kill \$RULEMANCER_PID' to stop it."



