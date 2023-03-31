# Best Effort Broadcast (BEB) Example

This example shows how to use the DISSE library to simulate a Best Effort Broadcast (BEB) module which sends
received packets to all other nodes in the network.

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

In order to run the tests, run the following command in the `beb` directory:

```bash
go test . -run TestBeb -v -l beb.log -u /dev/null
```

> Note: The `-u /dev/null` flag is used to disable the UML sequence diagram generation.

## Generating the UML Sequence Diagram

In order to generate the UML sequence diagram, you need to download the PlantUML jar file from [here](http://plantuml.com/download) and a Java runtime environment (JRE) from [here](http://openjdk.java.net/install/).

Then, create a .env file in the `beb` directory and add following lines:

```bash
DISSE_JAVA_PATH=<path to java executable>
DISSE_PLANTUML_JAR=<path to plantuml.jar>
```

A png file will be automatically generated in the `beb` directory after running the simulation from the uml file.
