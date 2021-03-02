# Golang-Challenge
Will use this readme to inform of some of the decisions taken.

## The Map
The TransparentCache uses a map in which is writes and reads since we are using many goroutines to access it we have to make it threadsafe, there were many options here

* Using a normal Map and a read/write MUTEX
* Use the sync.Map
* Make a map of mutexes to have a mutex for each key

I've chose the first one since I think it's the simplest and actually performs quite well and maintains type safety, if however this service were to be run on a host with more than 4 cores available then the contention for the locks will be detrimental to the performance and a sync.Map should be used. The map of mutexes just felt like overkill here where with luck most of the operations will be reads.
## Test
Use to test and search for race conditions.
* go test ./... -v  to test
* go test ./... -race 

I've only expanded one test with more elements to check the order o the GetPricesFor and added one to see that errors are returned for the same function since that one was missing. Examples and Benchamark tests should be added as well in the future. The benchmark test can be done by making sum (of the elapsed time) on a variable when waiting, that way the tests do not need to have the timeout and can be run faster.