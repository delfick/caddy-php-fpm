Caddy PHP-FPM 
=============

A module to run and supervise php-fpm in the background

This module is optimised for running a PHP app using caddy in Google Cloud Run
in that it will not let Caddy start until php-fpm has started. Otherwise Google
Cloud run will see your port is open and start serving to it.

How it works
------------

This will spawn a php-fpm process in the background and prevent caddy from starting till it's
up and running. When Caddy stops, the php-fpm process will also be stopped.

Full HTTP Caddyfile example
--------------------------

Caddyfile::

    {
        # Must be in global options
        php-fpm {
            cmd php-fpm -y fpm.conf
            sock_location path/to/fpm.sock
            start_timeout 10s
        }
    }

    mysite.com {
        @trailing-slash {
            path_regexp dir (.+)/$
        }
        rewrite @trailing-slash {re.dir.1}

        root * /var/www

        try_files {path} {path}.php {path}/index.php =404
        php_fastcgi unix/path/to/fpm.sock
    }

Building it
-----------

Use the `xcaddy` tool to build a version of caddy with this module::

    xcaddy build \
        --with github.com/delfick/caddy-php-fpm

Credits
-------

This package is a stripped down version of https://github.com/baldinof/caddy-supervisor,
which itself is a continuation of https://github.com/lucaslorentz/caddy-supervisor
which only supports Caddy v1.

Thanks @lucaslorentz and @baldinof for the original works!
