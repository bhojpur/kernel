// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

process.chdir('/tmp');

var StdOutFixture = require('./fixture-stdout.js');
var stdOutFixture = new StdOutFixture();

//set up log server
var _log = [];
stdOutFixture.capture( function onWrite (string, encoding, fd) {
  _log.push(string);
  return true;
});
var StdErrFixture = require('./fixture-stderr');
var stdErrFixture = new StdErrFixture();
stdErrFixture.capture( function onWrite (string, encoding, fd) {
  _log.push(string);
  return true;
});

const PORT=9967;
var http = require('http');
function serveLogs(request, response){
    response.end(_log.join(""));
}
var server = http.createServer(serveLogs);
server.listen(PORT, function(){
    console.log("Log server started on: http://localhost:%s", PORT);
});

console.log("Bhojpur Kernel v0.0 boostrapping beginning udp broadcast...");
tring = require('querystring');
  var options = {
    hostname: "169.254.169.254",
    path: '/latest/user-data',
    method: 'GET',
  };
  var req = http.request(options, function(res) {
    console.log('Status: ' + res.statusCode);
    console.log('Headers: ' + JSON.stringify(res.headers));
    res.setEncoding('utf8');
    res.on('data', function (body) {
      console.log('Response: ' + body);
      env = JSON.parse(body);
      Object.keys(env).forEach(function(key) {
        var val = env[key];
        process.env[key] = val;
        console.log("Set env var: "+key+"="+val)
      });
      console.log("Bhojpur Kernel v0.0 boostrapping finished!\ncalling main");
      //CALL_NODE_MAIN_HERE
    });
  });
  req.end();