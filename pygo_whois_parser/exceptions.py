class WhoisParserError(Exception):
    """Base class for all exceptions raised by WhoisParser."""
    pass

class InvalidInputError(WhoisParserError):
    """Raised when the input to the WHOIS parser is invalid."""
    pass

class WhoisParsingError(WhoisParserError):
    """Raised when the WHOIS parsing library returns an invalid result."""
    pass

class JSONParsingError(WhoisParserError):
    """Raised when the JSON returned by the WHOIS parser is invalid."""
    pass
