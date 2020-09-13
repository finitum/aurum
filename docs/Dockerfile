# Apidoc
FROM node:slim as node

WORKDIR /doc
COPY core core
COPY webclient webclient
#COPY apidoc.json

RUN yarn global add apidoc typedoc
RUN apidoc -i core/web/ -o /out/api

# Js doc
WORKDIR webclient
RUN yarn install
RUN typedoc webclient/ts --out /out/webclient


# Go doc
FROM golang
RUN go get golang.org/x/tools/cmd/godoc
COPY core /doc/core
WORKDIR /doc/core
COPY docs/generate_godocs.sh .
RUN ./generate_godocs.sh
COPY --from=node /out/ /usr/share/nginx/html/

# Final image
FROM nginx
COPY docs/index.html /usr/share/nginx/html/
COPY --from=golang /doc/core/godoc /usr/share/nginx/html/
