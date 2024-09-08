from setuptools import setup, find_packages

with open("README.md", "r") as f:
    long_description = f.read()

setup(
    name="pygo_whois_parser",
    version="0.0.1",
    description="A Python WHOIS parser leveraging a Go-based shared library for efficient domain data extraction.",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/nsb-coder/pygo-whois-parser",
    author="Nikola Stankovic",
    author_email="nikola.stankovic28991@gmail.com",
    license="MIT",
    classifiers=[
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3.10",
        "Operating System :: OS Independent",
    ],
    package_dir={"": "app"},
    packages=find_packages(where="app"),
    install_requires=[],
    package_data={
        "pygo_whois_parser": ["app/pygo_whois_parser/src/go-whois-parser/go-whois-parser.so"]
    },
)