# Code for the QUIC overview talk
## Requires
1. Basic go 
1. Basic VSCode -- another IDE will work, but the .vscode, which contains a way to debug with launch parameters, will be less helpful

### Certificates
QUIC and http3 only work over https. The [cert creation script](certs/createcerts.sh) makes creating the scripts easier.
Run the script from the [certs subdir](certs) Don't use these certs for anything else, and its not a bad idea to keep the duration relatively short. (Seven Days).

### Usage
go run . [opts]
#### options
| option | description | default |
| -| - | -: |
| -s | Server to Run |0|
| -c | Client to Run -1 means don't run a client | -1 |

Not all server/client pairs work with each other.
Server 1 and 2 work with a browser

Server 0 is a basic HTTPS server
Server 1 is a server that serves HTTP/3
