FROM iron/base

WORKDIR /app

# Add the binary
ADD go-chat /app/
ENTRYPOINT ["./go-chat"]
