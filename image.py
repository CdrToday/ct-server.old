import os
import tempfile
import uuid
import base64
import toml
from pathlib import Path
from flask_cors import CORS
from flask import Flask, escape, request, jsonify
from qiniu import Auth, put_file, etag, config

CONF = toml.load('./config.toml')['qiniu']

# utils
def cdn(pic):
    # config
    access_key = CONF['ak']
    secret_key = CONF['sk']
    bucket_name = CONF['bucket']

    q = Auth(access_key, secret_key)
    key = uuid.uuid4()
    token = q.upload_token(bucket_name, key, 3600)

    # save data
    png = base64.b64decode(pic)
    name = str(uuid.uuid4())
    path = str(Path.home()) + '/tmp/cache'

    if os.path.exists(path) == False :
        os.makedirs(path)

    path = path + name + '.png'
    with open(path, 'wb') as f:
        f.write(png)

    # upload
    ret, info = put_file(token, key, path)
    os.remove(path)
    
    if info.status_code == 200:
        return name

    return 'err'

# mainp
app = Flask(__name__)
CORS(app)

@app.route('/upload', methods=['GET', 'POST'])
def upload():
    j = request.json
    image = j['Image']
    res = cdn(image)
    
    return jsonify(
        msg=res
    )

app.run(host='0.0.0.0', port=7070)
