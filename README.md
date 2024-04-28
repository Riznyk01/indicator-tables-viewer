To run the program, run launcher.exe; it checks for updates and automatically updates the program if needed.  
The launcher must be running with the CFG_PATH environment variable or with the CFG_PATH flag (the flag takes precedence over the environment variable).
Flag or ENV is required.

If the program is updated by the launcher or if an update is not necessary, the launcher runs viewer.exe with  
the CONFIG_PATH environment variable, which equals CFG_PATH flag or environment variable specified at launcher startup.

Thanks to the code_path variable in the config, the launcher can be run from the IDE with full functionality in RUN  
mode (without biuld). It can run viewer.exe at the end of its work from the code_path because in RUN mode in  
the IDE, the program runs from the IDE TEMP directory which doesn't contain viewer.exe and resources.

The LocalPath config variable is needed if the program works in local mode (database on the localhost in the  rpogram dir).
In this case, the program determines from which directory it is run and uses the database and resources from its directory. Firebird must be installed in this case on the localhost.