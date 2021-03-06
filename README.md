Simple Cache Server
===================

For Go Learning Challenge - Simple Cache Server <http://www.topcoder.com/challenge-details/30046225/>

Simple Server Launch
--------------------

A simple server can be run with the following command: -

> go run scs.go -port {server port:11212} -items {maximum no. items to store:65535}

In addition to the commands listed in the spec at the above URL, additional commands can be added by
supplying the scs.RegisterCommands option to the scs.NewStore.

Testing
-------

A test file is supplied to test the default commands. To run the test, load the server with the following: -

> go run scs.go -port {port}

and then run the test: -

> go run test.go -port {port}

Due to the test checking stats, the test will only work on a fresh server with the default item limit.
