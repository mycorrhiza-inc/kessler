FROM python:3.12


ENV VIRTUAL_ENV=/usr/local
# use UV
# ADD --chmod=755 https://github.com/astral-sh/uv/releases/latest/download/uv-installer.sh /install.sh
# RUN /install.sh && rm /install.sh

WORKDIR /app
RUN apt-get update && apt-get install -y pandoc
COPY requirements.txt /app/requirements.txt
# Uncomment for debugging python depenancy issue
# COPY requirements-unlocked.txt /app/requirements.txt

RUN pip install --no-cache -r requirements.txt 
COPY . . 

# expose to docker network
EXPOSE 5555

# Command should be in dockerfile
CMD ["uvicorn","main:app", "--host", "0.0.0.0", "5555", "--workers", "4"]
