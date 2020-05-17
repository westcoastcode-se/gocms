# Example

This is a very simple example on the structure of how the CMS works.

The example shows how:
* You configure the authorization engine
* Configure ACL (roles give access to urls)
* Users (password is base64 encoded)
* Urls cached when running in public mode
* Search for specific page types and list them dynamically (news)

## Start

### Certificate 

Generate public and private keys used by the authorization engine. 
The config directory contains a pre-generated certificate. Please don't use that
in your own environment. 

You can generate a certificate with:
 
```bash
openssl req -x509 -newkey rsa:4096 -nodes -out cert.pem -keyout key.pem -days 365 && \
    openssl rsa -in key.pem -pubout > key.pub
```

### Start

Start the main.go and make sure to set the "example" folder as the home directory

