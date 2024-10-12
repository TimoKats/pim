![logo](https://github.com/TimoKats/pim/blob/main/.github/logo.png)

![example workflow](https://github.com/timokats/pim/actions/workflows/test.yaml/badge.svg)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-red.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![GitHub tag](https://img.shields.io/github/tag/TimoKats/pim?include_prereleases=&sort=semver&color=cyan)](https://github.com/TimoKats/pim/releases/)
[![stars - pim](https://img.shields.io/github/stars/TimoKats/pim?style=social)](https://github.com/TimoKats/pim)
[![forks - pim](https://img.shields.io/github/forks/TimoKats/pim?style=social)](https://github.com/TimoKats/pim) 

Pim (which stands for Process IMprover) is a task orchestrator that adds functionality to cron. And yes, it's named after a certain text editor I used to write the tool. It's written in Go and tested on Linux, although most funcionalities should also work on Windows.

## Getting started
You can install pim with `go install github.com/TimoKats/pim` (assuming you have go installed and GOPATH set correctly). Currently, pim is not available in any package repositories. Next, you can setup your tasks in `~/.pim/process.yaml`.

## Usage

### process.yaml
All your tasks can be defined in your process file, which sits in `~/.pim/process.yaml`. Here you can find a template of this file with some explanation of the possible attributes. Note, new ideas are welcome!

```yaml
max_logs: 50 # optional: max logs pim will store. Trims based on FIFO logic. Defaults to none.
only_store_errors: true # optional: if true, only logs of erronious runs are stored. Defaults to false
process:
  - run:
    name: fetch-repo # name of your run, mandetory
    schedule: '@times;8:00;10:00;16:00;18:00;21:55' # schedules prefixed with @times will run every day at the selected time.
    directory: /home/user/code/ # optional: set a directory to run the code in.
    catchup: true # optional. if true, pim will do a catchup run on startup if the computer was off when last scheduled.
    command: git fetch # the actual command you want to run.

  - run:
    name: scraper
    schedule: '@start+20' # schedules prefixed with @start run on startup. +20 means wait 20 seconds after startup to run.
    duration: 60 # optional: max number of seconds the command can run. After which, program is gracefully exited
    command: python3 scraper.py

  - run:
    name: change-wallpaper
    schedule: '*/5 * * * *' # You can also have some good-ol cron strings :)
    command: ./change-wallpaper.sh
```
