# Basic Port Scanner

This is a very basic TCP port scanner written in Go as an exercise to learn Go.
It happens to be fairly efficient at finding open TCP ports because it scans
in order of frequency (according to the nmap-services list), and because it
uses the nice concurrency features built into Go.

## TODO

- [ ] Allow scanning specified ports using same syntax as nmap
- [ ] Allow scanning multiple hosts using same syntax as nmap
- [ ] Show progress bar?