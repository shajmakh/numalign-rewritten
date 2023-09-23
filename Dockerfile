# this is a multi-stage dockerfile; run tests and build binary; generate the image.

# Use go 1.20 as the base image 
FROM golang:1.20


# create the working directory in the containercfor the next commands;
# in the next commands we address it as the current directory: ./
WORKDIR /app

COPY . .
# to be able to run make target from inside the "dockerfile" like this, 
# the project content (mainly go.mod and go.sum) must be placed in non-root directory,
# thus set the working directoty to some new directory called app 
RUN make tests
RUN make 

# Copy the current directory contents into the container 
COPY /app/build/numalign numalign

ENTRYPOINT ["/bin/sh"]

