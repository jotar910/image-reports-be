#!/bin/sh

# Populate databases.
sh scripts/seed-users.sh;
sh scripts/seed-reporter.sh;
sh scripts/seed-processing.sh;