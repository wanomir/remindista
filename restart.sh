#!/bin/bash

task docker:down
rm -rf db_data/
task run
