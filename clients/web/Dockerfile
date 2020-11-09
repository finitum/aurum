FROM node:13-slim as builder

WORKDIR /webclient
# Install dependencies
COPY package.json .
COPY yarn.lock .
RUN yarn install
COPY . .

# build
RUN yarn parcel build index.html

FROM nginx
EXPOSE 80

COPY --from=builder /clients/web /usr/share/nginx/html
