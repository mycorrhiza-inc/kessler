FROM python:3.12-bookworm

ENV VIRTUAL_ENV=/usr/local
# use UV
ADD --chmod=755 https://github.com/astral-sh/uv/releases/latest/download/uv-installer.sh /install.sh
RUN /install.sh && rm /install.sh
 
COPY requirements.txt /requirements.txt
 
RUN /root/.cargo/bin/uv pip install --no-cache -r requirements.txt 

# expose to docker network
EXPOSE 5505

CMD ["uvicorn","main:app", "--host", "0.0.0.0", "5505", "--workers", "4"]
