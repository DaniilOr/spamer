from ml import predict
from flask import Flask, request, jsonify
from flask import make_response
app = Flask(__name__)
@app.route('/', methods=['POST'])
def pred():
    if request.method == 'POST':
        url = request.get_json().get("url")
        print(url)
        #print(model.predict_proba(url))
        #print(url)
        res = predict(url)
        if res == "ham":
            res = "ok"
        return  make_response(jsonify({'Verdict':res}), 200)

if __name__ == "__main__":
    app.run(host="0.0.0.0")
