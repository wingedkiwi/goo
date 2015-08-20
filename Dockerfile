FROM golang
COPY goo/ /go/src/goo/
COPY authenticator/ /go/src/authenticator/

RUN go get -d -v authenticator && go install -v authenticator
RUN go get -d -v goo && go install -v goo

EXPOSE 80

CMD ["goo"]
