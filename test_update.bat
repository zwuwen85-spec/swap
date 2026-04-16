@echo off
cd /d "D:\毕业设计\CampusSwapShop\backend"
echo Testing user update endpoint...
curl -X PUT http://localhost:8080/api/v1/user/info ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer YOUR_TOKEN_HERE" ^
  -d "{\"nickname\":\"test\",\"qq\":\"123\",\"wechat\":\"test456\"}"
echo.
pause
