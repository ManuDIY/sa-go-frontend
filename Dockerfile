FROM alpine
COPY ./index.html /
COPY ./sa-go-frontend /
CMD ["/sa-go-frontend"]
