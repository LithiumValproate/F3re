# Use an official Python runtime as a parent image
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Copy the dependencies file to the working directory
COPY requirements.txt .

# Install any needed packages specified in requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Copy the application source code
COPY ./src /app

# Make port 5000 available to the world outside this container
EXPOSE 5000

# Define environment variable
ENV NAME F3reAuth

# Run main.py when the container launches
# Uvicorn is a lightning-fast ASGI server, running on host 0.0.0.0 makes it accessible from outside the container
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5000"]
