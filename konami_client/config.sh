#!/bin/sh

gobgp global rib add 10.0.0.0/8 community 65000:38  # ↑
gobgp global rib del 10.0.0.0/8
gobgp global rib add 10.0.0.0/8 community 65000:38  # ↑
gobgp global rib add 10.0.0.0/8 community 65000:40  # ↓
gobgp global rib del 10.0.0.0/8
gobgp global rib add 10.0.0.0/8 community 65000:40  # ↓
gobgp global rib add 10.0.0.0/8 community 65000:37  # ←
gobgp global rib add 10.0.0.0/8 community 65000:39  # →
gobgp global rib add 10.0.0.0/8 community 65000:37  # ←
gobgp global rib add 10.0.0.0/8 community 65000:39  # →
gobgp global rib add 10.0.0.0/8 community 65000:66  # B
gobgp global rib add 10.0.0.0/8 community 65000:65  # A
gobgp global rib del 10.0.0.0/8

gobgp neighbor 192.168.0.64 reset
