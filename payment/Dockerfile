FROM python:3.10-alpine3.15

WORKDIR /code
COPY . .

ENV FLASK_APP=app.py
ENV FLASK_RUN_HOST=0.0.0.0

RUN apk add --no-cache gcc musl-dev linux-headers

RUN pip install --upgrade pip
RUN pip install -r requirements.txt

EXPOSE 4242

CMD ["flask", "run", "--port=4242"]