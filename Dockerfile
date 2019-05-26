FROM alpine
COPY ./index.html /
COPY ./sa-go-frontend /
EXPOSE 8090
CMD ["/sa-go-frontend"]
