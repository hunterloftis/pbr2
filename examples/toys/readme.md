From the package root (..):

```
$ heroku create
$ heroku stack:set container
$ heroku labs:enable runtime-dyno-metadata
$ git push heroku master
$ heroku scale web=1:performance-l worker=10:performance-l
$ heroku open
```