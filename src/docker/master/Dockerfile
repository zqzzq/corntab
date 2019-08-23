FROM istionightly/base_debug

COPY job-master /usr/local/bin/
COPY master.json /etc/
ENTRYPOINT ["/usr/local/bin/job-master", "--config","/etc/master.json"]
