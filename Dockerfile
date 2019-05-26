FROM scratch
COPY index.html /
COPY sa-go-grontend /
CMD [/sa-go-grontend]

