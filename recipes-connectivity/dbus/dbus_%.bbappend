# override default so it starts first, before init-ifupdown
INITSCRIPT_PARAMS = "start 01 5 3 2 . stop 20 0 1 6 ."
