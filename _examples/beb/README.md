# Best Effort Broadcast (BEB) Example

This example shows how to use the DISSE library to simulate a Best Effort Broadcast (BEB) module which sends
received packets to all other nodes in the network, including itself.

## Specification

With best-effort broadcast, the burden of ensuring reliability is only on the sender.
Therefore, the remaining processes do not have to be concerned with enforcingthe reliability of received messages.
On the other hand, no delivery guarantees are offered in case the sender fails.

### Module properties

- **Validity:** If a correct node broadcasts a message, then all correct nodes eventually handle that message.
- **No duplication:** No message is handled more than once by a correct node.
- **No creation:** If a correct node handles a message, then that message was previously broadcast by some correct node.

## Running Simulation

In order to run the simulation, you need to have Go installed on your machine. You can download Go from [here](https://golang.org/dl/).

Then, build the simulation by running the following command in the `beb` directory:

```bash
go build -o beb.ds
```

Finally, run the simulation by running the following command in the `beb` directory:

```bash
./beb.ds -l beb.log -u beb.uml
```

For usage information, run the following command:

```bash
./beb.ds -h
```

## Running Tests

In order to run the all tests, run the following command in the `beb` directory:

```bash
go test . -v -l /dev/null -u /dev/null
```

> Note: The `-l /dev/null -u /dev/null` flags disable the UML sequence diagram generation and logging.

To run a specific test, run the following command in the `beb` directory:

```bash
go test . -run <test name> -v -l beb.log -u beb.uml
```

For example, to run the `TestBeb/TestValidity` test, run the following command:

```bash
go test . -run TestBeb/TestValidity -v -l beb.log -u beb.uml
```

## Generating the UML Sequence Diagram

In order to generate the UML sequence diagram, you need to download the PlantUML jar file from [here](http://plantuml.com/download) and a Java runtime environment (JRE) from [here](http://openjdk.java.net/install/).

Then, create a .env file in the `beb` directory and add following lines:

```bash
DISSE_JAVA_PATH=<path to java executable>
DISSE_PLANTUML_JAR=<path to plantuml.jar>
```

A png file will be automatically generated in the `beb` directory after running the simulation from the uml file.
