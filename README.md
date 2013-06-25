nagios-receiver
===============

Nagios-receiver is a daemon to receive nagios/check-mk results pushed through https/http.

It is designed to be placed behind an nginx or apache reverse proxy.
The reverse proxy can authenticate the clients with either HTTP auth or client
certificates and should write the username/CN into a HTTP header.
The header for authentication is configurable, and defaults to `X-REMOTE-USER`.

nagios-receiver writes the POST data from a client into a file which is constructed from
`{spool_dir}/{file_prefix}HTTP_HEADER`. (eg. /var/lib/icinga/ramdisk/nagios-receiver.myhost)
This makes it essentialy maintenance free as username/password or certificate
handling is entierly done by the reverse proxy.

**WARNING: nagios-receiver IS NOT INTENDED TO BE USED WITHOUT REVERSE PROXY!**


Overview
--------

What you need:

* Installed nagios-receiver daemon on the server
* Working reverse proxy to nagios-receiver for authentication (client cert or HTTP auth)
* Check-mk installation with added nagios-receiver.mk to `/etc/check_mk/conf.d/`
* A few hosts tagged with `nagpush` in check-mk
* Those hosts pushing the output of `check_mk_agent` to your server


Install
-------

Section on how to compile nagios-receiver and put it on the server

Install go on your system:

    apt-get install golang-go
    export GOPATH=~/go
    mkdir ~/go

Downloading and compile the lastest version of nagios-receiver

    go get github.com/gebi/nagios-receiver
    rsync -cvt $GOPATH/bin/nagios-receiver root@SERVER:/srv/nagios-receiver


Apache config
-------------

Sample config for apache used as a reverse proxy for nagios-receiver

Enable needed modules

    a2enmod headers

Config snippet

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

Examples on what is needed on the client side to push check-mk output to the server

Add http authentication information to ~/.netrc

    machine SERVER login USER password PASS

Simple crontab script

    `* * * * * check_mk_agent | curl -s --netrc --data-binar @- https://SERVER/check-receiver/`

Command to submit check information with proxy

    export https_proxy="http://USER:PASS@proxy-url:proxy-port"
    check_mk_agent | curl -s --netrc --data-binar @- https://SERVER/check-receiver/


Nagios/Icinga Config
--------------------

The ramdisk is shared between nagios/icinga and nagios-receiver which is also in the group nagios.
To be able to write to the ramdisk.

From /etc/rc.local:

    mountpoint -q /var/lib/icinga/ramdisk || mount -t tmpfs tmpfs /var/lib/icinga/ramdisk -o uid=104,gid=108,mode=2771,size=250m

On Debian
    uid = uid of nagios
    gid = gid of nagios


Daemon With Runit
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

    rsync -cvt $GOPATH/src/github.com/gebi/nagios-receiver/nagios-receiver.runit root@SERVER:/etc/sv/nagios-receiver/run


Debug with
---------

Some informations on how to debug nagios-receiver if you want to hack on it.

    # start daemon with default.conf and debug output
    ./nagios-receiver -debug

    # write data for user "foo"
    /bin/echo -e 'a\nb\nc\nd' |curl --data-binary @- -u foo:bar -H "X-REMOTE-USER: foo" http://localhost:8443

