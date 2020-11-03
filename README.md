iRApi
=====


iRacing Golang API.

**NOTICE** 

The API for iRacing is not officially supported and using it outside the web app is not supported.
The API may break at any time thus breaking this library and any integrations using it.

Maybe one day iRacing will provide an officially supported API for things like session results
and driver profiles.


Features:
---------

* Support for request tracing and context
* Overridable HTTP client/transport
* Pluggable sources for credentials
* Callbacks to inspect/modify requests and responses


Examples:
---------

See the [`examples/`](./examples) folder for full examples.