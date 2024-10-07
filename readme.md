# Compiled0

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [License](#license)

## Introduction

The Compiler Service is a robust and efficient solution for compiling various programming languages. It provides a RESTful API that allows developers to compile code remotely, making it ideal for integration into IDEs, online coding platforms, and educational tools.

## Features

- Support for multiple programming languages (e.g., C++, Java, Python, JavaScript)
- RESTful API for easy integration
- Secure sandboxed compilation environment
- Customizable compilation options
- Detailed error reporting and output capture
- Scalable architecture for handling multiple compilation requests

## Getting Started

### Prerequisites

- Docker (version 20.10 or later)
- Docker Compose (version 1.29 or later)
- Git

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/Ammyy9908/compiler-service.git
   cd compiler-service
   ```

2. Build and start the service using Docker Compose:
   ```
   docker-compose up --build
   ```

The service will be available at `http://localhost:8080`.

## Usage

To compile code, send a POST request to the `/compile` endpoint with the following JSON payload:

```json
{
    "code": "cHJpbnQoImhlbGxsbyIp",  // base64 encoded code
    "language_id": 1,  // language id
    "request_id": "114ecba7-61fb-4ae8-ad15-f67b44c07da7",  // request id
    "stdin": ""  // stdin
}
```

To get the result of the compilation, send a GET request to the `/submissions/<request_id>` endpoint.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
