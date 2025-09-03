# JobTrail

A lightweight job tracker built in Go, designed to handle incoming data from a Firefox extension, store it in an SQL database, and gracefully manage shutdowns. On receiving termination signals (e.g., Ctrl+C), the application automatically exports all database records to a CSV file and then truncates the table to reset the state. 
