A simple go server for single commands.

Usage:   
`$ go run server.go PORT 'COMMAND'`    
 \- variables are indicated with {NAME}

Example usage:   
`$ go run server.go 9900 'echo Someone said <{string}>'`    
 \- the server is then called from URL http://localhost:9900?string=Hello%20world

Output example:   
`$ curl http://localhost:9900?string='Hello%20world'`    
`result: Someone said <Hello world>`



### Security notice

Use at your own risk.. Be aware of injection risks for certain commands.


### TODO
* Call with client's local files (send as POST request)
* Piping?
