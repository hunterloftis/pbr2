```
$ heroku create pbr-toys
$ git push heroku master
$ heroku scale web=1:performance-l worker=10:performance-l
$ heroku open
```