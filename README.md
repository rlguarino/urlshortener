# URL Shortener
> Ross Guarino (rssguar@gmail.com)

An experiment in complicating things by turning a perfectly good implementation into Micro-Services.

## Services

### Frontend Service

The Frontend service handles all incoming user requests. It is the only service that
serves user requests directly.

#### Endpoints
    * / Accepts: GET; Displays Create new Shortened URL Form.
    * /new Accepts: POST; Creates new Shortened URL.
    * /r/{key} Accepts: GET; Directs the client to the target address or 404
    * /s/{key} Accepts: GET; Displays statistics for the shortened URL.

### URL Service

The URL service maps shortened URLs to the targets and is responsible for
generating the shortened version of the URLs.

#### Endpoints
    * /v1/route/{key} Accepts: GET; Returns `Redirect` JSON object for the short url.
    * /v1/route/ Accepts: POST; Create a new short URl for the Target URL specified in the `Reidrect` in the request body.
    * /info Accepts: GET; Display info about the Redis connection.
    * /env Accepts: GET; Display info about the environment.

### Stats Service

The Stats Service records clicks and generates statics for each URL based on the clicks.
Currently the Stats service only counts the number of clicks but it records more
information so it can be expanded in the future.

#### Endpoints
    * /v1/click Accepts: POST; Record a `Click` object.
    * /v1/stats/{key} Accepts: GET; Returns Statistics for the specified shortened URL.
    * /env Accepts: GET; Display info about the environment.


## Example Configuration & Kubernetes Deployment
Each service has an example configuration which will be used to start the service
in local development mode if not overridden. To specify a new config use the `-c`
or `-config` command line arguments.

There is an example Kubernetes deployment in ./kubernetes-example/.
Look at the README in that directory for more information.

## TODO 
    * Better Document the APIs
    * Add CSS to the Frontend
    * Cache target pages with a CacheService
    * More detailed statistics
    * Users
