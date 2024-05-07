## Table of Contents

- [Screenshots](#screenshots)
- [Running the program](#running-the-program)

## Screenshots

![Screenshot 1](screenshots/screenshot1.png)
![Screenshot 2](screenshots/screenshot2.png)
![Screenshot 3](screenshots/screenshot3.png)

## Running the Program

To run the program:

1. Ensure that the CFG_PATH environment variable is set or specify the CFG_PATH flag.
   CFG_PATH is a full path to the config file in .toml format. 
2. Execute launcher.exe.  

The launcher automatically checks for updates and updates the program if necessary.
Autoupdate is enabled if 'auto_update' is set to true in the config.toml file. 
If the program is updated or no update is required, the launcher executes viewer.exe using the CONFIG_PATH
environment variable, which is equal to the CFG_PATH flag or environment variable specified at launcher
startup.
The code_path variable in the configuration allows the launcher to be run from the IDE with full
functionality in RUN mode (without building). It can execute viewer.exe from the code_path because
in RUN mode in the IDE, the program runs from the IDE TEMP directory, which doesn't contain viewer.exe
and resources.  

The LocalPath config variable is necessary if the program operates in local mode (database on the
localhost in the program directory). In this scenario, the program determines the directory from
which it is run and uses the database and resources from its directory. Firebird must be installed
on the localhost in this case.

