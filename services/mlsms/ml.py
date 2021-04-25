import numpy as np
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.metrics import accuracy_score
import pandas as pd
import pickle
model = pickle.load(open("naive.pkl", "rb"))
import string
import nltk
from nltk.corpus import stopwords
from sklearn.model_selection import train_test_split
from google_trans_new import google_translator
import  requests
nltk.download('stopwords')
def process(message):
    nopunc = [char for char in message if char not in string.punctuation]
    nopunc = ''.join(nopunc)

    # remove any stopwords
    return [word for word in nopunc.split() if word.lower() not in stopwords.words('english')]

messages = pd.read_csv('spam.csv',encoding="ISO-8859-1")
messages = messages.drop(['Unnamed: 2', 'Unnamed: 3', 'Unnamed: 4'], axis=1)
messages = messages.rename(columns={'v1':'label', 'v2':'message'})
X_train,X_test, y_train, y_test = train_test_split(messages['message'], messages['label'], test_size = 0.3, random_state = 0)
cv = CountVectorizer(analyzer=process)
X_train = cv.fit_transform(X_train)



def predict(sms):
    translator = google_translator()
    sms = translator.translate(sms)
    transformed = cv.transform([sms])
    pred = model.predict(transformed)

    return pred[0]

