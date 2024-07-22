FROM python:3.12.4-bookworm


ENV VIRTUAL_ENV=/usr/local
# use UV
# ADD --chmod=755 https://github.com/astral-sh/uv/releases/latest/download/uv-installer.sh /install.sh
# RUN /install.sh && rm /install.sh

WORKDIR /app
RUN apt-get update 
RUN apt-get install -y pandoc=2.17.1.1-2~deb12u1 --no-install-recommends
COPY requirements.txt /app/requirements.txt
# Uncomment for debugging python depenancy issue
# COPY requirements-unlocked.txt /app/requirements.txt

RUN pip install --no-cache -r requirements.txt 
COPY . . 

# expose to docker network
EXPOSE 5555

# Command should be in dockerfile
CMD ["uvicorn","main:app", "--host", "0.0.0.0", "5555", "--workers", "4"]
