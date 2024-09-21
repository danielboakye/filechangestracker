repeat
	display dialog "Manage the 'filechangestracker' service:" buttons {"Start Service", "Stop Service", "Exit"} default button "Start Service"
	set userChoice to button returned of the result

	if userChoice is "Start Service" then
		try
			-- Change to the target directory and run 'go build'
			do shell script "cd ~/workspace/filechangestracker && /usr/local/bin/go build"
			
			-- Ask for the user's password
			display dialog "Enter your password to start the service:" default answer "" with hidden answer
			set userPassword to text returned of the result
			
			-- Run the built binary in the background with 'sudo'
            do shell script "cd ~/workspace/filechangestracker && echo " & quoted form of userPassword & " | sudo -S nohup ./filechangestracker > /dev/null 2>&1 & disown"

			
			display dialog "Service started!" buttons {"OK"} default button "OK"
		on error errMsg
			display dialog "Error starting service: " & errMsg buttons {"OK"} default button "OK"
		end try
	else if userChoice is "Stop Service" then
		try
			-- Stop the service
			display dialog "Enter your password to stop the service:" default answer "" with hidden answer
			set userPassword to text returned of the result
			
			-- Using `sudo` to stop the background process
			do shell script "echo " & quoted form of userPassword & " | sudo -S pkill -f ./filechangestracker"
			
			display dialog "Service stopped!" buttons {"OK"} default button "OK"
		on error errMsg
			display dialog "Error stopping service: " & errMsg buttons {"OK"} default button "OK"
		end try
	else if userChoice is "Exit" then
		exit repeat
	end if
end repeat
