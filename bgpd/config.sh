#!/bin/sh

gobgp global rib add 0.0.0.0/0

gobgp policy prefix add default 0.0.0.0/0
gobgp policy statement default-accept add condition prefix default
gobgp policy statement default-accept add action accept
gobgp policy add default-accept default-accept

gobgp neighbor add 192.168.0.71 as 65000  # Update accordingly here
gobgp global policy export add default-accept default reject

# this is deadly slow
gobgp mrt inject global $*
