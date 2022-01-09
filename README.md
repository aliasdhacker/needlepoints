# needlepoints

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

 

