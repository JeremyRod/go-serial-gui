# go-serial-gui
A Serial GUI to help users choose, setup and monitor serial ports. 

## Use Case 

This gui is a simple implementation of a serial console/monitor for COM port and eventually VCOM ports. 
For more complicated console and terminal usage, another terminal application will be a better use case, however for a simple connection, logging and console output, this application will do nicely. 

## How to use

Start the gui and configure the port using the appropriate parameters.
This will spawn a new window with a console where the UART data will be printed.
The field below is used to send data to the connected device.

The File button on the main window will allow the user to create a file in a location of their choosing to log all read data from the port.
