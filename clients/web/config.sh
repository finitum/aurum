#!/bin/sh

# With help of https://www.brandonbarnett.io/blog/2018/05/accessing-environment-variables-from-a-webpack-bundle-in-a-docker-container/
echo "export default { API_URL: '${API_URL}' };" > ts/Config.ts

cat ts/Config.ts
exec nginx -g 'daemon off;'
