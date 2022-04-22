Deployment
==========

This is a stripped down version of what I actually use.

General idea is:

1. building/make_caddy - Create a docker image that has the caddy binary we use
2. building/make_bundle - Makes a folder containing everything we need to
   deterministically create a deployment. So that we may make the final docker
   image for staging and production at different times and get the same result
3. building/make_final - Creates the final image to deploy
4. ci/deploy - Puts the image in google cloud run

Locally you can run ``./build_image`` to make the same image.

And you may run ``./run`` to run the caddy and php-fpm without needing docker
at all.
