# -*- coding: utf-8 -*-
# @Author: BigHead
# 自动根据钉钉用户（webhook消息中的senderStaffId字段）id@人

import hmac
import hashlib
import base64
import requests
import urllib
import time
from lib import common_log


logger = common_log.get_logger(__name__)


class Ding:
    # appkey, appsecret： 钉钉开发者后台获取,可以查询access token
    # robot_token, robot_secret: 用于发送消息到客户端
    def __init__(self, userid, appkey: str, appsecret: str, robot_token: str, robot_secret: str):
        self.userid = userid
        self.app_key =  appkey
        self.app_secret = appsecret
        self.robot_token = robot_token
        self.robot_secret= robot_secret

    def get_access_token(self):
        access_token = requests.get(
            url='https://oapi.dingtalk.com/gettoken',
            params={
                'appkey': self.app_key,  # 应用信息中的AppKey
                'appsecret': self.app_secret
            }
        )

        a_c = access_token.json()['access_token']
        return a_c

    def user_info(self):
        url = "https://oapi.dingtalk.com/user/get"
        params = {
            "access_token": self.get_access_token(),
            "userid": self.userid
        }
        response = requests.get(url, params=params)
        logger.info(f"user:{self.userid}查询结果{response.text}")
        return response.json()

    def user_phone(self):
        mobile = self.user_info().get("mobile")
        if mobile:
            return True, mobile
        else:
            return False, ""


    def gen_sign(self):
        timestamp = int(time.time() * 1000)
        data = (str(timestamp) + '\n' + self.robot_secret).encode('utf-8')
        hmac_code = hmac.new(self.robot_secret.encode('utf-8'), data, digestmod=hashlib.sha256).digest()
        sign = urllib.parse.quote_plus(base64.b64encode(hmac_code))
        return timestamp,sign

    def ding_markdown(self, content):
        times, sign = self.gen_sign()
        url = f'https://oapi.dingtalk.com/robot/send?access_token={self.robot_token}&timestamp={times}&sign={sign}'
        ok, mobile = self.user_phone()
        if not ok:
            logger.info("提取手机号失败，放弃@指令")
            # return
        else:
            content = content + f"@{mobile}"
        data = {
            "msgtype": "markdown",
            "markdown": {
                "title": content,
                "text": content,
            },
            "at": {
                "atMobiles": [
                    mobile
                ],
                "atUserIds": [
                    self.userid
                ],
                "isAtAll": "false"
            }
        }
        logger.info(f"请求dingding发送消息，内容: {data}")
        r = requests.post(url, json=data)
        logger.info(f"发送markdown消息结果: {r}")

    def ding_text(self):
        pass


if __name__ == "__main__":
    userDing = Ding('16060611799966130', "xxx", "xxx-PcZrYJTwft", "xx","xx")
    mess = f'''
### 项目发布开始  
> 目标环境: test1,共12个项目
- k8s_dtk_php_www_user_center
- k8s_dtk_php_www_user_center4
'''
    userDing.ding_markdown(mess)

    # userDing.ding_markdown('''
## 发布成功
# ---
# - 项目: k8s_dtk_php_www_user_center
# - 项目: k8s_dtk_php_www_user_center2
## 发布失败
# ---
# - 项目: k8s_dtk_php_www_user_cente5： 代码无变化,停止发版8s_dtk-php-dy-goods-admin 分支: dtk-php-dy-goods-admin_2023SXSZD-1140
# - 项目: k8s_dtk_php_www_user_center2: jjj
#     ''')
