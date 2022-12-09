import stripe
import os
import pymongo
from flask import Flask, request, render_template

YOUR_DOMAIN = os.getenv("DOMAIN", "http://localhost:4242")

# This is your Stripe CLI webhook secret for testing your endpoint locally.
stripe.api_key = "sk_test_51L0tL5SJYv4HUS02q8sw8iaGzTrAifySQ3eluMsT3PHVCDJhZ7c0cD7h5nvcjgKKfonboUsP66zITCfJ890xxFhX00F4CGVgS7"

ENDPOINT_SECRET = (
    "whsec_62d9de44628dad668f972fa7c81aa939c722cb46a74cbd97e34ed8660073c903"
)
mongo_certificate = os.getenv("MONGODB_CERTIFICATE", "")

uri = (
    "mongodb+srv://hosteldb.e3ayhyn.mongodb.net/?retryWrites=true&w=majority&authSource=%24external&authMechanism=MONGODB-X509&tlsCertificateKeyFile="
    + mongo_certificate
)

mongo_client = pymongo.MongoClient(uri, serverSelectionTimeoutMS=5000)
students_collection = mongo_client.get_database("University").get_collection("Students")

app = Flask(__name__)


@app.route("/", methods=["GET"])
def get_form():
    return render_template("form_entry.html")


@app.route("/", methods=["POST"])
def get_payment_intent():
    form_id = request.form.get("form-id")
    student = students_collection.find_one(filter={"forms": {"$exists": form_id}})
    if not student:
        return "No such form"
    try:
        intent = stripe.PaymentIntent.retrieve(student["forms"][form_id]["payment_ref"])
    except Exception as ex:
        return str(ex)

    return render_template("checkout.html", client_secret=intent.client_secret)


@app.route("/success", methods=["GET"])
def success():
    return render_template("success.html")


@app.route("/webhook", methods=["POST"])
def webhook():
    payload = request.get_data()
    sig_header = request.headers.get("STRIPE_SIGNATURE", "")
    event = None

    try:
        event = stripe.Webhook.construct_event(payload, sig_header, ENDPOINT_SECRET)
    except ValueError:
        return "Invalid payload", 400
    except Exception:
        return "Invalid signature", 400

    event_dict = event.to_dict()
    if event_dict["type"] == "payment_intent.succeeded":
        intent = event_dict["data"]["object"]
        print("Succeeded: ", intent["id"])
        # Fulfill the customer's purchase
    elif event_dict["type"] == "payment_intent.payment_failed":
        intent = event_dict["data"]["object"]
        error_message = (
            intent["last_payment_error"]["message"]
            if intent.get("last_payment_error")
            else None
        )
        print("Failed: ", intent["id"]), error_message
        # Notify the customer that payment failed

    return "OK", 200
