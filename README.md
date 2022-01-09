# needlepoints

This application is used to track and spent points.  "Gifting" points is done by adding a new transaction.  Specify the amount of points, the payer, and the timestamp.

Use the "Spend" url to spend these points, the application spends the oldest points first, not letting any one payer points balance go below zero.

This is a go application, built on go 1.16.2.  To run the application checkout this source code, install golang 1.16.2, browse to this source directory and type "go build ."

This should create an executable called "needlepoint.exe"

I have included the EXE in the repo, you should be able to just run this application from the EXE.


3 endpoints exposed
 (Default listen, localhost port 4000)
 
GET Request http://localhost:4000/
 - Get all payer points balances

POST Request http://localhost:4000/newTxn
 - Post a new transaction, sample request
 - curl http://localhost:4000/newTxn    
   --include     \
   --header "Content-Type: application/json"   \
   --request "POST"    \
   --data '{"payer": "DANNON","points": -200, "timestamp": "2020-10-31T15:00:00Z"}'
 
POST Request http://localhost:4000/spend
 - Spend points, sample transaction
 - curl http://localhost:4000/spend    \
    --include     \
    --header "Content-Type: application/json"   \
    --request "POST"  \
    --data '{"points": 5000}'

 

