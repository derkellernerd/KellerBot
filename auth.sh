#!/usr/bin/env bash
scopes="user:read:chat user:write:chat user:bot"

auth_result=$(curl -v --location 'https://id.twitch.tv/oauth2/device' \
		   --form "client_id=\"${TWITCH_CLIENT_ID}\"" \
		   --form "scopes=\"${scopes}\"")
echo $auth_result

uri="${auth_result}" jq -r '.verification_uri'
device_code="${auth_result}" jq -r '.device_code'
echo "Open: ${uri}"
echo "Insert Device Code: ${device_code}"

read -p "Finish?" finish

result=$(curl --location 'https://id.twitch.tv/oauth2/token' \
     --form "client_id=\"${TWITCH_CLIENT_ID}\"" \
     --form "scopes=\"${scopes}\"" \
     --form "device_code=\"${device_code}\"" \
     --form 'grant_type="urn:ietf:params:oauth:grant-type:device_code"' | jq)
