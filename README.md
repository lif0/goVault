<div align="center">
<img src=".github/assets/goVault_poster_round.png" > 
<h2 align="center">✨Blazing-fast in-memory database written in Go✨</h2>

<p align="center">
<img src="https://github.com/lif0/goVault/actions/workflows/coverallsapp.yml/badge.svg"/>
<img src="https://goreportcard.com/badge/github.com/lif0/goVault" />
<a href='https://coveralls.io/github/lif0/goVault?branch=main'>
  <img src='https://coveralls.io/repos/github/lif0/goVault/badge.svg?branch=main' alt='Coverage Status' />
</a>
<a href="https://github.com/lif0/goVault/blob/main/LICENSE">
  <img src="https://img.shields.io/github/license/lif0/goVault.svg?logo=github&labelColor=rgb(64%2C70%2C78" />
</a>
</p>

<h3 align="center">Please leave a ⭐ as motivation if you liked the DB 😄
<br>🌪️Currently a WIP and in Active development.</h3>

<h>If you have feature request feel free to open an [Issue](https://github.com/lif0/goVault/issues/new/choose)</h4>
</div>
<br />


## About DB
**goVault**  is designed to provide quick and efficient access to data. Built with simplicity and performance in mind, that is perfect for applications that demand speed and low latency.

## Features

- **High Performance**: Optimized for rapid data storage and retrieval.
- **Simple API**: Easy-to-use commands for storing, retrieving, and deleting data.
- **In-Memory Storage**: All data is stored in memory, ensuring low-latency operations.
- **Lightweight**: Minimal dependencies, built entirely in Go.

## Getting Started

### Installation

You can pull the Docker image directly from Docker Hub and start using goVault immediately:

1. Pull the Docker image:
```
docker pull lif0/govault:latest
```

2. Create a `config.yaml` file with your custom configuration. Example:
```yaml
engine:
  max_threads: 4
network:
  address: ":7777"
logging:
  level: "debug"
  output: "./goVault/log/output.log"
  stdout: false
wal:
  flushing_batch_length: 100
  flushing_batch_timeout: "10ms"
  max_segment_size: "10MB"
  data_directory: "./goVault/data/wal"
```

3. Run the container:
```
docker run -p 8080:{port-from-config} lif0/govault:latest
```

4. Now your goVault instance is running and ready to use on ```localhost:{port-from-config}```.

### Usage

The grammar of the query language in the form of eBNF:
```
query = set_command | get_command | del_command

set_command = "SET" argument argument
get_command = "GET" argument
del_command = "DEL" argument
argument    = punctuation | letter | digit { punctuation | letter | digit }

punctuation = "\*" | "/" | "_" | ...
letter      = "a"  | ... | "z" | "A" | ... | "Z"
digit       = "0"  | ... | "9"
```

#### Basic Commands

- **SET**: Store a key-value pair.
  ```
  SET key value
  ```
- **GET**: Retrieve the value for a key.
  ```
  GET key
  ```
- **DEL**: Delete a key-value pair.
  ```
  DEL key
  ```  

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request to improve goVault.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

---

Made with ❤️ in Go.
