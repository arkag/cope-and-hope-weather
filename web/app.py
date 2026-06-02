import json
import boto3
from flask import Flask, render_template, request, jsonify
from mangum import Mangum
from a2wsgi import WSGIMiddleware

app = Flask(__name__)


@app.route("/")
def index():
    return render_template("index.html")


@app.route("/search")
def search():
    city = request.args.get("city")
    mode = request.args.get("mode", "cope")

    # Mock an API Gateway V1 Event for the Go backend
    event = {
        "path": "/weather",
        "httpMethod": "GET",
        "queryStringParameters": {
            "city": city,
            "mode": mode
        }
    }

    try:
        # Call the Go API directly over AWS internal network
        lambda_client = boto3.client('lambda', region_name='us-east-1')
        response = lambda_client.invoke(
            FunctionName="CopeAndHopeWeatherAPI",
            InvocationType="RequestResponse",
            Payload=json.dumps(event)
        )

        payload = json.loads(response['Payload'].read().decode('utf-8'))

        if payload.get("statusCode") != 200:
            err_msg = f"API returned {payload.get('statusCode')}: {payload.get('body')}"
            return jsonify({"error": err_msg}), payload.get("statusCode", 502)
        return jsonify(json.loads(payload["body"]))
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/health")
def health():
    return jsonify({"status": "ok"})


# Lambda entry point requires ASGI, but Flask is WSGI.
asgi_app = WSGIMiddleware(app)
handler = Mangum(asgi_app)

if __name__ == "__main__":
    app.run(debug=True, port=5000)
