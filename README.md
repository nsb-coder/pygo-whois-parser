# pygo-whois-parser
A Python library for fast WHOIS data parsing using a Go-based shared library under the hood.

For a project that required fast parsing of large amounts of WHOIS data, existing Python libraries were too slow, so I used Go to significantly boost performance.

## Features
- Fast WHOIS parsing using a Go-based shared library.
- Input validation for non-empty strings.
- Converts WHOIS data into a Python dictionary.

## Requirements
Python 3.x

## Notes
This project is still a work in progress.
Support for additional TLDs (Top-Level Domains) is coming soon.

## License
This project is licensed under the MIT License.

## Installation

```

```

## Usage

```
from whois_parser import WhoisParser

# Initialize the parser
parser = WhoisParser()

# Raw WHOIS data as a string
whois_raw_data = "Your raw WHOIS record here..."

# Parse the WHOIS data
parsed_data = parser.parse(whois_raw_data)
print(parsed_data)
```

### Example Output

```
{
  "domain_name": "example.com",
  "registrar": "Example Registrar",
  "creation_date": "2020-01-01",
  "expiration_date": "2025-01-01",
  ...
}
```