FROM python:3.10-alpine3.20
WORKDIR /main/
RUN apk update && apk add build-base libpq-dev
RUN pip3 install psycopg2 prometheus_client
COPY main.py .
ENTRYPOINT [ "python3", "main.py" ]
