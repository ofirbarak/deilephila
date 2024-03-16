# Tests
## Setup
We use docker-compose to setup our architecture.
Run: `./setup_env.sh` to run the MongoDB and Neo4j.

## Integration Tests
In the main folder of the project run: `make integration-test`.

## Benchmark Tests
We used the framework of Golang to benchmark the system. 
You can run the benchmark tests from the main folder of the project by running: `make benchmark-test`.
