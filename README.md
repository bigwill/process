# Using process

After building the project, you can hear process make noise by doing this:

$ sox <AUDIO FILE> -t f64 -r 48k -c 1 - | process | play -t f64 -r 48k -c 1 -

sox and play commands are available in the sox project (http://sox.sourceforge.net/)