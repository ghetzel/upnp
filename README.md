# upnp

Golang UPNP implementation to open gateway ports and get WAN address

Available functons include:

- AddMappingPort
- DelMappingPort
- ExternalIPAddress

## NOW WITH COMMAND LINE UTILITY !!! WOW !!!

```
make
./bin/upnpc map 9876/tcp
# <do internet>
./bin/upnpc unmap 9876/tcp
```
