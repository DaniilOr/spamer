import itertools
from sklearn.metrics import mean_squared_error,confusion_matrix, precision_score, recall_score, auc,roc_curve
from sklearn.model_selection import train_test_split
import pandas as pd
import numpy as np
import random
import math
from collections import Counter
from sklearn import metrics
import matplotlib.pyplot as plt
from sklearn.metrics import classification_report,confusion_matrix,accuracy_score
from lightgbm import LGBMClassifier
from urllib.parse import urlparse
from googlesearch import search
import re
from urllib.parse import urlparse
from tld import get_tld
import os.path
from sklearn.preprocessing import LabelEncoder

#Use of IP or not in domain
def having_ip_address(url):
    match = re.search(
        '(([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.'
        '([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\/)|'  # IPv4
        '((0x[0-9a-fA-F]{1,2})\\.(0x[0-9a-fA-F]{1,2})\\.(0x[0-9a-fA-F]{1,2})\\.(0x[0-9a-fA-F]{1,2})\\/)' # IPv4 in hexadecimal
        '(?:[a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}', url)  # Ipv6
    if match:
        # print match.group()
        return 1
    else:
        # print 'No matching pattern found'
        return 0


def abnormal_url(url):
    hostname = urlparse(url).hostname
    hostname = str(hostname)
    match = re.search(hostname, url)
    if match:
        # print match.group()
        return 1
    else:
        # print 'No matching pattern found'
        return 0
def google_index(url):
    site = search(url, 5)
    return 1 if site else 0
def no_of_dir(url):
    urldir = urlparse(url).path
    return urldir.count('/')
def no_of_embed(url):
    urldir = urlparse(url).path
    return urldir.count('//')
def shortening_service(url):
    match = re.search('bit\.ly|goo\.gl|shorte\.st|go2l\.ink|x\.co|ow\.ly|t\.co|tinyurl|tr\.im|is\.gd|cli\.gs|'
                      'yfrog\.com|migre\.me|ff\.im|tiny\.cc|url4\.eu|twit\.ac|su\.pr|twurl\.nl|snipurl\.com|'
                      'short\.to|BudURL\.com|ping\.fm|post\.ly|Just\.as|bkite\.com|snipr\.com|fic\.kr|loopt\.us|'
                      'doiop\.com|short\.ie|kl\.am|wp\.me|rubyurl\.com|om\.ly|to\.ly|bit\.do|t\.co|lnkd\.in|'
                      'db\.tt|qr\.ae|adf\.ly|goo\.gl|bitly\.com|cur\.lv|tinyurl\.com|ow\.ly|bit\.ly|ity\.im|'
                      'q\.gs|is\.gd|po\.st|bc\.vc|twitthis\.com|u\.to|j\.mp|buzurl\.com|cutt\.us|u\.bb|yourls\.org|'
                      'x\.co|prettylinkpro\.com|scrnch\.me|filoops\.info|vzturl\.com|qr\.net|1url\.com|tweez\.me|v\.gd|'
                      'tr\.im|link\.zip\.net',
                      url)
    if match:
        return 1
    else:
        return 0

def suspicious_words(url):
    match = re.search('PayPal|login|signin|bank|account|update|free|lucky|service|bonus|ebayisapi|webscr',
                      url)
    if match:
        return 1
    else:
        return 0

#First Directory Length
def fd_length(url):
    urlpath= urlparse(url).path
    try:
        return len(urlpath.split('/')[1])
    except:
        return 0
def tld_length(tld):
    try:
        return len(tld)
    except:
        return -1
def digit_count(url):
    digits = 0
    for i in url:
        if i.isnumeric():
            digits = digits + 1
    return digits
def letter_count(url):
    letters = 0
    for i in url:
        if i.isalpha():
            letters = letters + 1
    return letters


def process_url(url):
    df = pd.DataFrame({"url": [url]})
    df['use_of_ip'] = df['url'].apply(lambda i: having_ip_address(i))
    df['abnormal_url'] = df['url'].apply(lambda i: abnormal_url(i))
    df['google_index'] = df['url'].apply(lambda i: google_index(i))
    df['count.'] = df['url'].apply(lambda i: i.count('.'))
    df['count-www'] = df['url'].apply(lambda i: i.count('www'))
    df['count@'] = df['url'].apply(lambda i: i.count('@'))
    df['count_dir'] = df['url'].apply(lambda i: no_of_dir(i))
    df['count_embed_domian'] = df['url'].apply(lambda i: no_of_embed(i))
    df['short_url'] = df['url'].apply(lambda i: shortening_service(i))
    df['count-https'] = df['url'].apply(lambda i: i.count('https'))
    df['count-http'] = df['url'].apply(lambda i: i.count('http'))
    df['count%'] = df['url'].apply(lambda i: i.count('%'))
    df['count?'] = df['url'].apply(lambda i: i.count('?'))
    df['count-'] = df['url'].apply(lambda i: i.count('-'))
    df['count='] = df['url'].apply(lambda i: i.count('='))
    # Length of URL
    df['url_length'] = df['url'].apply(lambda i: len(str(i)))
    # Hostname Length
    df['hostname_length'] = df['url'].apply(lambda i: len(urlparse(i).netloc))
    df['sus_url'] = df['url'].apply(lambda i: suspicious_words(i))
    df['tld'] = df['url'].apply(lambda i: get_tld(i, fail_silently=True))
    df['tld_length'] = df['tld'].apply(lambda i: tld_length(i))
    df['count-digits'] = df['url'].apply(lambda i: digit_count(i))
    df['count-letters'] = df['url'].apply(lambda i: letter_count(i))
    df['fd_length'] = df['url'].apply(lambda i: fd_length(i))

    df = df.drop("tld", 1)
    return df[['use_of_ip', 'abnormal_url', 'google_index', 'count.', 'count-www', 'count@',
               'count_dir', 'count_embed_domian', 'short_url', 'count-https',
               'count-http', 'count%', 'count?', 'count-', 'count=', 'url_length',
               'hostname_length', 'sus_url', 'fd_length', 'tld_length', 'count-digits',
               'count-letters']]

