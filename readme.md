![logo](https://github.com/TimoKats/pim/blob/main/.github/logo.png)

## Abstract
![example workflow](https://github.com/timokats/pim/actions/workflows/test.yaml/badge.svg)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-red.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![GitHub tag](https://img.shields.io/github/tag/TimoKats/pim?include_prereleases=&sort=semver&color=cyan)](https://github.com/TimoKats/pim/releases/)
[![stars - pim](https://img.shields.io/github/stars/TimoKats/pim?style=social)](https://github.com/TimoKats/pim)
[![forks - pim](https://img.shields.io/github/forks/TimoKats/pim?style=social)](https://github.com/TimoKats/pim) 

Pim (which stands for Process IMprover) is a minimally invasive task orchestrator that runs on Windows and Linux. Besides regular cron scheduling, it adds additional features like:

- Run commands on startup (with optional delay).
- Run commands in specific directories.
- Run commands directly based on their set name (aliases).
- Doing catchup runs (if computer was turned off during scheduled run).
- Access to run logs.
- Etc...Ideas welcome!

## Getting started
You can install pim with `go install github.com/TimoKats/pim@latest` (assuming you have go installed and GOPATH set correctly). Currently, pim is not available in any package repositories. Next, you can setup your schedule. On linux systems, this will be in `~/.pim/process.yaml`. On windows, you can do the same in `C:\Users\<username>\.pim\process.yaml`

## Usage

### process.yaml
All your tasks can be defined in your process file, which sits in `~/.pim/process.yaml`. Here you can find a template of this file with some explanation of the possible attributes. In summary, there are some overall optional settings, and there is a list of "runs". Here, each run represents a command you want to schedule.

```yaml
max_logs: 50 # optional: max logs pim will store. Trims based on FIFO logic. Defaults to none.
only_store_errors: true # optional: if true, only logs of erronious runs are stored. Defaults to false
process:
  - run:
    name: fetch-repo # name of your run, mandetory
    schedule: '@times;8:00;10:00;16:00;18:00;21:55' # schedules prefixed with @times will run every day at the selected time.
    directory: /home/user/code/ # optional: set a directory to run the code in.
    catchup: true # optional. if true, pim will do a catchup run on startup if the computer was off when last scheduled.
    command: git fetch

  - run:
    name: scraper
    schedule: '@start+20' # schedules prefixed with @start run on startup. +20 means wait 20 seconds after startup to run.
    duration: 60 # optional: max number of seconds the command can run. After which, program is gracefully exited
    command: python3 scraper.py

  - run:
    name: change-wallpaper
    schedule: '*/5 * * * *' # You can also have some good-ol cron strings :)
    command: ./change-wallpaper.sh

  - run:
    name: python-program # Without cron schedule, you can run this on command
    directory: /home/user/code/python/
    command: python3 main.py -d some_parameter

```

### Commands
Pim is used from the command line. This is an overview of the commands/flags you can use. Note, the format is `pim <<command>>` and `pim --<<flag>>`. The gif below shows some pim commands being used in practice.

Commands:
- run <\<command-name\>>: Runs a command by the name defined in your process YAML.
- start: Starts the cron schedule defined in your process YAML. NOTE: Logs are sent to log files in ~/.pim/logs/ and not to standard output when using start.
- stop: Stops the cron schedule started by running: pim start.
- ls: Lists all the commands and their characteristics defined in your process YAML.
- log <\<optional:run-id\>>: Show all logs, or a log of a specific run.
- status: Show if pim is currently running and its pid.
- clean: Clean log files.

Flags:
- info/i: Outputs some information about this Pim installation.
- help/h: Well...if you see this message you probably typed this...
- version/v: Shows version of this Pim installation.
- license/l: Shows the license of this Pim installation.

![pim-example-optmized](https://github.com/user-attachments/assets/85d69829-739c-49e1-a75e-63c2ff3d2324)


### Running in background
If you want to run `pim start` automatically you have a couple of options. First, you can add `(pim start &)` to your .bashrc file. This will start a pim session whenever you open your terminal (note, it prevents running multiple times if you have multiple windows). Next, you can also use something like crontab to start on reboot. A third option would be to run pim through a systemd file. Finally, it's also possible to just run `pim start` and `pim stop` whenever you want it to be on. E.g. in servers that never shut down.

> For Windows users: Perhaps you can run pim as a service, but there's no way for me to test that given I don't have a Windows machine.

## FAQ
- Q: I have some scheduling requirements that pim doesn't support.
    - A: I would love to know! If you write code, feel free to make a PR. Else, open something on the discussion board.
- Q: Can I install pim without go?
    - A: Currently not. If pim actually becomes popular I will work on putting it in package managers.
- Q: What does "minimally invasive task orchestration" mean?
    - A: Programs like systemD control all your processes. Pim can be started/stopped on command and only impacts the processes you add to process.yaml.
- Q: What advantages does pim have over anacron?
    - A: Anacron can only run once a day (or less often)
- Q: What advantages does pim have over crontab?
    - A: Crontab can't do catchup runs if the computer was turned off, and only accepts cron strings.
 - Q: What advantages does pim have over SystemD?
    - A: If you already have systemD and are happy with it, maybe not so much. However, pim is cross platform and minimally invasive, so you can copy your process.yaml to any other OS and it should work. SystemD runs all your programs and is linux-oriented. Therefore, it's less flexible.
