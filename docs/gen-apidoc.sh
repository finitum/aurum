#!/bin/sh
~/.yarn/bin/apidoc -i ../core/web/ -o ./apidoc
xdg-open ./apidoc/index.html
