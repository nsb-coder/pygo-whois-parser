import ctypes
import json
import os
from typing import Dict
import platform


class WhoisParser:
    """Class for parsing WHOIS data using a shared Go-based library."""

    def __init__(self):
        """
        Initialize the WhoisParser by loading the shared library.
        """
        so_file_path = self._get_shared_library_path()
        self.lib = ctypes.CDLL(so_file_path)
        self.lib.ParseWhois.argtypes = [ctypes.c_char_p]
        self.lib.ParseWhois.restype = ctypes.c_char_p

    def _get_shared_library_path(self):
        """Determine the correct shared library based on the OS."""
        base_dir = f"{os.path.dirname(os.path.abspath(__file__))}/go-whois-parser/"
        system = platform.system()

        if system == "Linux":
            return os.path.join(base_dir, "go-whois-parser.so")
        elif system == "Darwin":
            return os.path.join(base_dir, "go-whois-parser.dylib")
        elif system == "Windows":
            return os.path.join(base_dir, "go-whois-parser.dll")
        else:
            raise OSError(f"Unsupported operating system: {system}")

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
