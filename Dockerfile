FROM ubuntu:latest

#setting the WORKDIR
WORKDIR /app

# Copy the mcluster-vcontroller binary to the working directory
COPY ./mcluster-vcontroller ./mcluster-vcontroller

# Ensure the binary is executable
RUN chmod +x ./mcluster-vcontroller

# Set the container to run the mcluster-vcontroller binary by default
ENTRYPOINT ["./mcluster-vcontroller"]

