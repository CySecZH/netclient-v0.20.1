#!/bin/bash
#chat.openai.com
#docker run -e "TOKENS=token1,token2,token3" your_image_name

#Define cleanup
cleanup() {
    nets=($(wg show interfaces))
    for net in "${nets[@]}"; do
        echo "deleting interface" "$net"
        ip link del "$net"
    done
}

#Trap SigTerm
trap 'cleanup' SIGTERM

echo "[netclient] joining network"

if [ -z "${SLEEP}" ]; then
    SLEEP=10
fi

# Check if the TOKENS environment variable is set
if [ -n "$TOKENS" ]; then
    IFS=',' read -ra TOKENS_ARRAY <<< "$TOKENS"
else
    echo "No tokens provided in the TOKENS environment variable. Exiting."
    exit 1
fi

for TOKEN in "${TOKENS_ARRAY[@]}"; do
    if [ -n "$TOKEN" ]; then
        TOKEN_CMD="-t $TOKEN"
    else
        TOKEN_CMD=""
    fi

    /root/netclient join $TOKEN_CMD 
    if [ $? -ne 0 ]; then
        echo "Failed to join with token: $TOKEN, quitting."
        exit 1
    fi
done

echo "[netclient] Starting netclient daemon"

/root/netclient daemon &

wait $!
echo "[netclient] exiting"
