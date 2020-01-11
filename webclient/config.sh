#!/bin/sh
# With help of https://www.brandonbarnett.io/blog/2018/05/accessing-environment-variables-from-a-webpack-bundle-in-a-docker-container/
echo "export default { API_URL: '${API_URL}' };" > js/config.ts
cat js/config.mjs
exec nginx -g 'daemon off;'
