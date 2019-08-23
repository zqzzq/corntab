FROM istionightly/base_debug

COPY job-worker /usr/local/bin/
COPY worker.json /etc/
ENTRYPOINT ["/usr/local/bin/job-worker", "--config","/etc/worker.json"]
