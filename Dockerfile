FROM iron/base

WORKDIR /app

# Add the binary
ADD bin/go-chat /app/
ENTRYPOINT ["./go-chat"]
