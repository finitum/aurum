FROM node:13-slim as builder

WORKDIR /build
# Install dependencies
COPY package.json .
COPY yarn.lock .
RUN yarn install
COPY . .

# build
RUN yarn parcel build index.html

FROM nginx
EXPOSE 80

COPY --from=builder /build/dist /usr/share/nginx/html
