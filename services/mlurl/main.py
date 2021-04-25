from preprocessing import process_url
from flask import Flask, request, jsonify
from flask import make_response
import pickle
MODEL_PATH = "lgboost.pk"
app = Flask(__name__)
model = pickle.load(open(MODEL_PATH, 'rb'))
def match(code):
    if code == 0:
        return "ok"
    if code == 1:
        return "defacement"
    if code == 2:
        return "phishing"
    if code == 3:
        return "ok"

@app.route('/', methods=['POST'])
def pred():
    if request.method == 'POST':
        url = request.get_json().get("url")
       # if url[:5] != "http":
       #     url = "http://" + url
        if url.split(".")[0][-3:] != "www":
        	url = url.split("://")[0] + "://www." + "://".join(url.split("://")[1:])
        print(url)
        url = process_url(url)
        #print(model.predict_proba(url))
        #print(url)
        return  make_response(jsonify({'Verdict':match(model.predict(url)) }), 200)

if __name__ == "__main__":
    app.run(host="0.0.0.0")

