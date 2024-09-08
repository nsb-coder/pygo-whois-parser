import ctypes
import json
from typing import Dict


class WhoisParser:
    """Class for parsing WHOIS data using a shared Go-based library."""

    def __init__(self):
        """
        Initialize the WhoisParser by loading the shared library.
        """
        self.lib = ctypes.CDLL("./go-whois-parser/go-whois-parser.so")
        self.lib.ParseWhois.argtypes = [ctypes.c_char_p]
        self.lib.ParseWhois.restype = ctypes.c_char_p

    def validate_input(self, whois_raw: str) -> None:
        """
        Validate the input for the WHOIS parser.

        Args:
            whois_raw (str): The raw WHOIS record as a string.

        Raises:
            ValueError: If the input is not a non-empty string.
        """
        if not isinstance(whois_raw, str):
            raise ValueError("whois_raw must be a string.")

        if not whois_raw.strip():
            raise ValueError("whois_raw cannot be an empty string.")

    def parse(self, whois_raw: str) -> Dict:
        """
        Parse a raw WHOIS record and return a dictionary of parsed data.

        Args:
            whois_raw (str): The raw WHOIS record as a string.

        Returns:
            Dict: A dictionary containing the parsed WHOIS data.

        Raises:
            ValueError: If the input is invalid.
            json.JSONDecodeError: If the returned JSON data is invalid.
            RuntimeError: If the parsing library returns an invalid result.
        """
        self.validate_input(whois_raw)

        result = self.lib.ParseWhois(whois_raw.encode("utf-8"))
        if not result:
            raise RuntimeError("Failed to parse WHOIS data.")

        parsed_json = ctypes.c_char_p(result).value.decode("utf-8")
        return json.loads(parsed_json)


parser = WhoisParser()
print(parser.parse(""))