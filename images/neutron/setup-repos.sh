#!/bin/bash
# Copyright (c) 2020 VEXXHOST, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -xe

apt-get install -y gnupg2

cat <<EOF | apt-key add -
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v2.0.22 (GNU/Linux)

mQINBF3u81cBEACbfsspk7WNkcXCn3N5T9VKYt/dmvSsEW8nIIf/iwV7dSISmruz
1b7bviqfekEvf37yiFwHVFxxS70ry/ofXp51X7RUVytrJY/hNMvr7C7zyNqM928+
c8TP3FjGsPvFiWw/L2JgGl9/4+OYW5yF3HabMOa63xbFPAU891o9HIN5YfFDZZWD
VNsMyXCUjVB9wy7anF77moqXuews1OmvMSArE7erLjAnC5HHGdTeZO7KfDCylqPB
oBWF3pzNU1Vu6wEq9vL5NDYglsbN7jmDA+8mS0SyAnFxvTsqjisR8gpNtPaatQqM
wTdOydscSoXS9MfpCrxPne0dBmpAlcVdI4hq1T4l9Osf2x5s+Kb9JxF+Q4V87n4q
8fjusePRIMxO7aZjFUEvL8uIzg7VvF3b1X9UXkS6LH2YPLOqOf3lhvyk5RwwMfHp
p99KOVrTWbaBYVKuxR17oWkYBPOPp+4ld8F6zSk36GK+lzPP8814X28kS357lg1y
4kla/CfNav3AXdnsZkCvJhrwwR8HCXwTYaF2TzrZPqv5TZB1k9iBuL2X52BSxobR
PvTTM00iZhipC/EsA7vQu4FOla/ySb/R6cfFIiDyOrDiOJ3+zlWDQ0uBikCP4lIY
uUB+uVIWd8F7Us1voqsqUrVL1CSu1cYn+NOhf12eZsA740wgUZfCU2qmGwARAQAB
tBhyZXBvIDxyZXBvQHZleHhob3N0Lm5ldD6JAjkEEwEIACMFAl3u81cCGy8HCwkI
BwMCAQYVCAIJCgsEFgIDAQIeAQIXgAAKCRDETupUFYbjJ7WxD/9HcMd9HMwg7WC3
eKSFeHGJXtN/0IuCJ6r3q10/dhb8QqqZ+Rnlr5CH4DAHdkhnL5+OvnHVYu/LVejX
17dZUS0uB+JXZpMdfsv8i2g/c8uxi2KPsRa3pxXudb+WhjbxhRxeMpsQNbMc5M5+
cYseUYj1nzTioDn9MQH43GcYBuhydiWsp7zRs2CNWrWJgwTOwnd/g4YV+9VWqshM
x+/N0bdD+LIT0MmYYGBaK6vBnM2kG6gcwc0ZMwMYHJk+MotuFNM7KDu06XWkp/Uq
8uzi7tZKHTa/kc+LrJrIOwLIkFH1uMvRZXma+JwASbcEW97YCUw/vhLa3AbZvCum
9QLHv28zyUXfo9QLEhkOGC/ykkYOSt0u/lznokpf840tmYHBCLavFzOPJ0Nc2T7Y
tCyEA5sV2UVI4hdBtwG1Vz8rAggDu0NWDW3BGyP0X2x1jddzzNRhevqQqcAe83Ei
XOOP1aunhtUKUe+sXLFOY0d3OK0RysKAn9kdxcZ9qqZdrKhj+dwuvMBeZPau0ZGT
t81b/zv6hiwA1b1b4X6EKz/aZwQyQ3/UUovM0KC9rMSzm5kKYWwcfkSDY6aLZtgc
GBc+auY+9Mwcp4V5kEH6zMXF4baJzMj2m7LFYlLRVofY5kxlrr86TAK0jMmiDOx8
AcjcZTiXBPNU8sK+VbsXvtB0Mel7Vw==
=hpXM
-----END PGP PUBLIC KEY BLOCK-----
EOF

cat <<EOF | tee /etc/apt/sources.list.d/vexxhost.list
deb http://repo.vexxhost.net/ buster main
EOF
