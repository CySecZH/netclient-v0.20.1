#!/bin/bash
# https://app.askcodi.com/chat/
#This script assumes that the docker environment variables for the tokens are named as TOKEN_1, TOKEN_2, etc. 
#You can adjust the grep "^TOKEN_*" pattern to match your specific variable naming convention.
# Define cleanup
cleanup() {
    nets=($(wg show interfaces))
    for net in "${nets[@]}"; do
        echo "Deleting interface $net"
        ip link del "$net"
    done
}

# Trap SigTerm
trap 'cleanup' SIGTERM

echo "[netclient] Joining network"

if [[ -z "$SLEEP" ]]; then
    SLEEP=10
fi

TOKENS=()
for token_var in $(env | grep -i "^TOKEN_*" | cut -d= -f1); do
    token_value="${!token_var}"
    if [[ -n "$token_value" ]]; then
        TOKENS+=("$token_value")
    fi
done

if [[ ${#TOKENS[@]} -eq 0 ]]; then
    echo "No token variables found. Exiting."
    exit 1
fi

echo "Found ${#TOKENS[@]} token(s): ${TOKENS[*]}"

for token in "${TOKENS[@]}"; do
    /root/netclient join -t "$token"
    if [[ $? -ne 0 ]]; then
        echo "Failed to join with token: $token. Quitting."
        exit 1
    fi
done

echo "[netclient] Starting netclient daemon"

/root/netclient daemon &

wait $!
echo "[netclient] Exiting"