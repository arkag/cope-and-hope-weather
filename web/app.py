import os
import requests
from flask import Flask, render_template, request, jsonify
from mangum import Mangum

app = Flask(__name__)

API_ENDPOINT = os.environ.get("API_ENDPOINT", "http://localhost:8080")


@app.route("/")
def index():
    return render_template("index.html")


@app.route("/search")
def search():
    city = request.args.get("city", "")
    mode = request.args.get("mode", "cope")

    if not city:
        return jsonify({"error": "city is required"}), 400

    try:
        resp = requests.get(
            f"{API_ENDPOINT}/weather",
            params={"city": city, "mode": mode},
            timeout=10,
        )
        resp.raise_for_status()
        return jsonify(resp.json())
    except requests.RequestException as e:
        return jsonify({"error": str(e)}), 502


@app.route("/health")
def health():
    return jsonify({"status": "ok"})


# Lambda entry point
handler = Mangum(app)

if __name__ == "__main__":
    app.run(debug=True, port=5000)
