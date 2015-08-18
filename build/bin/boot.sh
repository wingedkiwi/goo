#!/bin/bash

set -aeup pipefail
if [ -f /etc/gopkg/env ]; then
  source /etc/gopkg/env
fi
gopkg -addr 0.0.0.0:80
