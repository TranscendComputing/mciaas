# MCIaaS

A (not yet production ready, rather a work in progress) REST api server interface to the [Packer](https://github.com/mitchellh/packer) application.

## Installation

This application is a Go program. Generally, assure that you have Go setup properly. Then clone this repository into $GOPATH/src/github.com/transcendcomputing/mciaas and type make. In a brief time, the mciaas executable will appear.

### Starting Up

To run the server, just execute with a root path (the directory in which Packer will execute) and a port on which to listen. For example:
    
    ./mciaas -r /home/myid/packer -p 8080

Thereafter the server listens for REST requests on that port.

## Usage

The REST API documentation is coming ... hang in, shortly we'll add it. The routes are as follows:

GET     /templates
GET     /templates/builders
GET     /templates/builders/:type
GET     /templates/provisioners
GET     /templates/provisioners/:type
DELETE  /packer/:user/:docId
GET     /packer/:user/:docId
POST    /packer/:user/:docId
PUT     /packer/:user
DELETE  /image/:user/:docId
GET     /image/:user/:docId
POST    /image/:user/:docId
PUT     /image/:user/:docId

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
