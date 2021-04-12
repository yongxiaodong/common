import hmac
import hashlib
import base64
import requests
import urllib
import time


def ding(title, options: dict):
    content = '\n'.join([f"+ **{k}**：{v}" for k, v in options.items()])
    secret = 'SEC789e87081fc800363b6520cd4d0dea9aeb151adf3ca22a83ab1ff67be8e408ba'
    token = '00c14a94c66e0223d71145811827f81468de45d22475b8f3c7e21c62c2c729b6'
    timestamp = int(time.time() * 1000)
    data = (str(timestamp) + '\n' + secret).encode('utf-8')
    hmac_code = hmac.new(secret.encode('utf-8'), data, digestmod=hashlib.sha256).digest()
    sign = urllib.parse.quote_plus(base64.b64encode(hmac_code))
    url = f'https://oapi.dingtalk.com/robot/send?access_token={token}&timestamp={timestamp}&sign={sign}'
    r = requests.post(url, json={
        "msgtype": "markdown",
        "markdown": {
            "title": title,
            "text": content,
        }
    })

if __name__ == "__main__":
    ding("测试标题", {"message": "内容"})