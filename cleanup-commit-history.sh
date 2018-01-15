#!/bin/bash
git checkout --orphan latest_branch
git add -A
git ci -am "initial commit"
git branch -D master
git branch -m master
git push -f origin master
git push -u origin master
