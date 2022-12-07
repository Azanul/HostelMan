# app.py
#
# Use this sample code to handle webhook events in your integration.
#
# 1) Paste this code into a new file (app.py)
#
# 2) Install dependencies
#   pip3 install flask
#   pip3 install stripe
#
# 3) Run the server on http://localhost:4242
#   python3 -m flask run --port=4242

import stripe

from flask import Flask, jsonify, request

# This is your Stripe CLI webhook secret for testing your endpoint locally.
endpoint_secret = (
    "whsec_62d9de44628dad668f972fa7c81aa939c722cb46a74cbd97e34ed8660073c903"
)

app = Flask(__name__)


@app.route("/webhook", methods=["POST"])
def webhook():
    event = None
    payload = request.data
    sig_header = request.headers["STRIPE_SIGNATURE"]

    try:
        event = stripe.Webhook.construct_event(payload, sig_header, endpoint_secret)
    except ValueError as ex:
        # Invalid payload
        raise ex
    except Exception as ex:
        # Invalid signature
        raise ex

    # Handle the event
    if event["type"] == "payment_intent.succeeded":
        payment_intent = event["data"]["object"]
    # ... handle other event types
    else:
        print(f"Unhandled event type {event['type']}")

    return jsonify(success=True)
