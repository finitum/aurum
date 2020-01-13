#!/bin/sh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
COMMON_PASSWORDS=`cat $DIR/../passwords/commonpasswordlist`

# With help of https://www.brandonbarnett.io/blog/2018/05/accessing-environment-variables-from-a-webpack-bundle-in-a-docker-container/
echo "export default { API_URL: '${API_URL}', COMMON_PASSWORDS:[$COMMON_PASSWORDS] };" > ts/config.ts
cat ts/config.ts
exec nginx -g 'daemon off;'
