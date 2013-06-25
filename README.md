nagios-receiver
===============

Daemon to receive nagios/check-mk results pushed through https/http

nagios-receiver is written to be placed behind an nginx or apache reverse proxy.
The reverse proxy can authenticate the clients with either http auth or client
certificates and should write the username/CN into an http header.
The header for authentication is configurable, default config file uses `REMOTE_USER`.

nagios-receiver writes the POST data from client into a file which consists of
`{spool_dir}/{file_prefix}HTTP_HEADER`.
This makes it essentialy maintenance free as username/password or certificate
handling is entierly done by the reverse proxy.

WARNING:
nagios-receiver IS NOT INTENDED TO BE USED WITHOUT REVERSE PROXY!


Apache config
-------------

a2enmod headers

    ProxyPass /check-receiver/ http://localhost:8443/
    <Location /check-receiver/>
        ProxyPassReverse /

        AuthName "Check Receiver"
        AuthType Basic
        AuthUserFile auth/check-receiver.auth
        Require valid-user

        RequestHeader set X-REMOTE-USER %{REMOTE_USER}s

        # limit POST to 1MB
        LimitRequestBody 1048576
    </Location>


Debug with
---------

    # start daemon with default.conf
    ./nagios-receiver -debug

    # write data for user "foo"
    /bin/echo -e 'a\nb\nc\nd' |curl --data-binary @- -u foo:bar -H "X-REMOTE-USER: foo" http://localhost:8443
