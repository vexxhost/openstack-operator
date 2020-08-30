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

cat <<EOF | apt-key add -
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1

mQINBFX4hgkBEADLqn6O+UFp+ZuwccNldwvh5PzEwKUPlXKPLjQfXlQRig1flpCH
E0HJ5wgGlCtYd3Ol9f9+qU24kDNzfbs5bud58BeE7zFaZ4s0JMOMuVm7p8JhsvkU
C/Lo/7NFh25e4kgJpjvnwua7c2YrA44ggRb1QT19ueOZLK5wCQ1mR+0GdrcHRCLr
7Sdw1d7aLxMT+5nvqfzsmbDullsWOD6RnMdcqhOxZZvpay8OeuK+yb8FVQ4sOIzB
FiNi5cNOFFHg+8dZQoDrK3BpwNxYdGHsYIwU9u6DWWqXybBnB9jd2pve9PlzQUbO
eHEa4Z+jPqxY829f4ldaql7ig8e6BaInTfs2wPnHJ+606g2UH86QUmrVAjVzlLCm
nqoGymoAPGA4ObHu9X3kO8viMBId9FzooVqR8a9En7ZE0Dm9O7puzXR7A1f5sHoz
JdYHnr32I+B8iOixhDUtxIY4GA8biGATNaPd8XR2Ca1hPuZRVuIiGG9HDqUEtXhV
fY5qjTjaThIVKtYgEkWMT+Wet3DPPiWT3ftNOE907e6EWEBCHgsEuuZnAbku1GgD
LBH4/a/yo9bNvGZKRaTUM/1TXhM5XgVKjd07B4cChgKypAVHvef3HKfCG2U/DkyA
LjteHt/V807MtSlQyYaXUTGtDCrQPSlMK5TjmqUnDwy6Qdq8dtWN3DtBWQARAQAB
tCpDZXBoLmNvbSAocmVsZWFzZSBrZXkpIDxzZWN1cml0eUBjZXBoLmNvbT6JAjgE
EwECACIFAlX4hgkCGwMGCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEOhKwsBG
DzmUXdIQAI8YPcZMBWdv489q8CzxlfRIRZ3Gv/G/8CH+EOExcmkVZ89mVHngCdAP
DOYCl8twWXC1lwJuLDBtkUOHXNuR5+Jcl5zFOUyldq1Hv8u03vjnGT7lLJkJoqpG
l9QD8nBqRvBU7EM+CU7kP8+09b+088pULil+8x46PwgXkvOQwfVKSOr740Q4J4nm
/nUOyTNtToYntmt2fAVWDTIuyPpAqA6jcqSOC7Xoz9cYxkVWnYMLBUySXmSS0uxl
3p+wK0lMG0my/gb+alke5PAQjcE5dtXYzCn+8Lj0uSfCk8Gy0ZOK2oiUjaCGYN6D
u72qDRFBnR3jaoFqi03bGBIMnglGuAPyBZiI7LJgzuT9xumjKTJW3kN4YJxMNYu1
FzmIyFZpyvZ7930vB2UpCOiIaRdZiX4Z6ZN2frD3a/vBxBNqiNh/BO+Dex+PDfI4
TqwF8zlcjt4XZ2teQ8nNMR/D8oiYTUW8hwR4laEmDy7ASxe0p5aijmUApWq5UTsF
+s/QbwugccU0iR5orksM5u9MZH4J/mFGKzOltfGXNLYI6D5Mtwrnyi0BsF5eY0u6
vkdivtdqrq2DXY+ftuqLOQ7b+t1RctbcMHGPptlxFuN9ufP5TiTWSpfqDwmHCLsT
k2vFiMwcHdLpQ1IH8ORVRgPPsiBnBOJ/kIiXG2SxPUTjjEGOVgeA
=/Tod
-----END PGP PUBLIC KEY BLOCK-----
EOF

cat <<EOF | tee /etc/apt/sources.list.d/ceph.list
deb https://download.ceph.com/debian-octopus/ buster main
EOF