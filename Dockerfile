FROM golang:1.19.3

ADD . /go/src

WORKDIR /go/src

CMD go run main.go

# This document is Licensed under Creative Commons CC0.
# To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights
# to this document to the public domain worldwide.
# This document is distributed without any warranty.
# You should have received a copy of the CC0 Public Domain Dedication along with this document.
# If not, see https://creativecommons.org/publicdomain/zero/1.0/legalcode.
