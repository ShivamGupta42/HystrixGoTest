#Installing Vegeta

````
brew update && brew install vegeta
````

#Using Vegeta
````
echo "GET http://localhost:5001/hysterix/test" | vegeta attack -duration=60s -rate=50 | vegeta report --type=text
````
