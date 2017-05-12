FROM scratch
WORKDIR /app
COPY main /app/
ENTRYPOINT ["./main"]