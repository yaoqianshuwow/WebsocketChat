import urllib.request, json

base = 'http://localhost:8888/api/v1'
users = [
    ('admin','123456','管理员'),
    ('u1','111111','用户一'),
    ('u2','111111','用户二'),
]
print("=== 注册 ===")
for name, pwd, nick in users:
    data = json.dumps({'username':name,'password':pwd,'nickname':nick}).encode('utf-8')
    req = urllib.request.Request(f'{base}/register', data=data, headers={'Content-Type':'application/json'})
    resp = json.loads(urllib.request.urlopen(req).read())
    print(f'{name}: code={resp["code"]} {resp.get("message","")}')

print("=== 登录 ===")
for name, pwd, _ in users:
    data = json.dumps({'username':name,'password':pwd}).encode('utf-8')
    req = urllib.request.Request(f'{base}/login', data=data, headers={'Content-Type':'application/json'})
    resp = json.loads(urllib.request.urlopen(req).read())
    print(f'{name}: id={resp.get("user_id")} nick={resp.get("nickname")}')
