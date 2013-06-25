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


Overview
--------

What you need:

    * Installed nagios-receiver daemon on the server
    * Working reverse proxy to nagios-receiver for authentication (client cert or http auth)
    * Check-mk installation with added nagios-receiver.mk to `/etc/check_mk/conf.d/`
    * A few hosts tagged with `nagpush` in check-mk
    * Those hosts pushing the output of `check_mk_agent` to your server


Install
-------

Installing go on your system:

    apt-get install golang-go
    export GOPATH=~/go
    mkdir ~/go

Downloading and compiling the lastest version of nagios-receiver

    go get github.com/gebi/nagios-receiver
    rsync -cvt $GOPATH/bin/nagios-receiver root@SERVER:/srv/nagios-receiver


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


Client Command
--------------

Add http authentication information to ~/.netrc

    machine SERVER login USER password PASS

Simple crontab script

    `* * * * * check_mk_agent | curl -s --netrc --data-binar @- https://SERVER/check-receiver/`

Command to submit check information with proxy

    export https_proxy="http://USER:PASS@proxy-url:proxy-port"
    check_mk_agent | curl -s --netrc --data-binar @- https://SERVER/check-receiver/


Daemon with runit
-----------------

On server:

    apt-get install runit

Create required user:

    # daemon user
    addgroup --system nagrecv
    adduser --system --home /nonexistent --no-create-home --disabled-login --ingroup nagrecv nagrecv
    adduser nagrecv nagios

    # logging user
    addgroup --system nagrecvlog
    adduser --system --home /nonexistent --no-create-home --disabled-login --ingroup nagrecvlog nagrecvlog

Configure runit and log service

    wget -O /usr/local/sbin/sva https://github.com/gebi/runit-toolkit/raw/master/usr_sbin/sva
    chmod 755 /usr/local/sbin/sva

    # configure runit service
    mkdir -p /etc/sv/nagios-receiver
    cd /etc/sv/nagios-receiver
    ln -s /var/run/sv.nagios-receiver supervise

    # configure logging service
    mkdir /etc/sv/nagios-receiver/log
    cd /etc/sv/nagios-receiver/log
    ln -s /var/run/sv.nagios-receiver.log supervise
    wget -O run https://github.com/gebi/runit-toolkit/raw/master/lib/scripts/common-log
    chmod 755 run
    mkdir conf
    echo nagrecvlog >conf/LOGUSER

From source:

    rsync -cvt nagios-receiver.runit root@SERVER:/etc/sv/nagios-receiver/run


Debug with
---------

    # start daemon with default.conf
    ./nagios-receiver -debug

    # write data for user "foo"
    /bin/echo -e 'a\nb\nc\nd' |curl --data-binary @- -u foo:bar -H "X-REMOTE-USER: foo" http://localhost:8443
