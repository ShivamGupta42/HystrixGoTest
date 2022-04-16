#Installing Vegeta

````
brew update && brew install vegeta
````

#Using Vegeta
````
echo "GET http://localhost:5001/hysterix/test" | vegeta attack -duration=60s -rate=50 | vegeta report --type=text
````

#Postman Collection
https://www.getpostman.com/collections/e4a44bbfe4af5f4ef354