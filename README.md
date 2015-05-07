## Getting Started
This site is based on revel, a web framework for the go language.  Deployment and setup operations are partially automated using rake (i.e. ruby make).  Go and rake are prerequisites for deploying the site.

### Rake commands:

    rake install: Installs dependencies and compiles tool for initializing human ppi predictions database. 

    rake launch_server[procs]: Sets the environmental variable which dictates how many system threads Go can utilize (procs) and launches the app. 
    
### Deployment Setup
Ubuntu comes with a service called Upstart, which isused for automatically starting services on system startup.  The configuration file for this web app is located on the host at: /etc/init/go-pcg-http.conf


The app uses Apache as a reverse proxy for handling requests.  

### Description of Contents

The default directory structure of a generated Revel application:

    myapp               App root
      app               App sources
        controllers     App controllers
          init.go       Interceptor registration
        models          App domain models
        routes          Reverse routes (generated code)
        views           Templates
      conf              Configuration files
        app.conf        Main configuration file
        routes          Routes definition
      db                Contains database file
      messages          Message files
      pgq_initdb        Source for program used to populate db
      public            Public assets
        css             CSS files
        js              Javascript files
        images          Image files
        predictions     Human PPI prediction archives
      tests             Test suites

app

    The app directory contains the source code and templates for your application.

conf

    The conf directory contains the application’s configuration files. There are two main configuration files:

    * app.conf, the main configuration file for the application, which contains standard configuration parameters
    * routes, the routes definition file.


messages

    The messages directory contains all localized message files.
