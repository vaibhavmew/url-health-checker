Problem Statement

We have n number of urls. All of them serving the same data. However the cost of an api call to each service is different.
All the urls are stored in an array with increasing order of cost.
By default we make the request to the first url. However if it goes down, we switch to the second url and so on till the last url post which 
we get internal server error. Once these apis go down we start a goroutine for each one of these and start calling their health endpoints.
The moment any of these endpoints are up again, we switch back to that url. 
The code gives us the ability yo have maximum uptime while minimizing the cost of api calls.
