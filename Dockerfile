FROM quay.io/brianredbeard/corebox

EXPOSE 8000
CMD ["/bin/resources", "serve", "-a", "0.0.0.0:8000"]

ADD resources /bin/resources
