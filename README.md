<p align="right">
 <img src="https://github.com/hellgate75/go-cron/workflows/Go/badge.svg?branch=master"></img>
&nbsp;&nbsp;<img src="https://api.travis-ci.com/hellgate75/go-cron.svg?branch=master" alt="trevis-ci" width="98" height="20" />&nbsp;&nbsp;<a href="https://travis-ci.com/hellgate75/go-cron">Check last build on Travis-CI</a>
 </p>

<p align="center">
<image width="150" height="147" src="images/clock.png"></image>&nbsp;
<image width="260" height="410" src="images/golang-logo.png">
&nbsp;<image width="150" height="121" src="images/terminal.png"></image>
</p><br/>
<br/>

# go-cron
Go Language scheduler

Scheduler that allow import and implementation of commands, it's more effective and it has more capabilities than the `unix cron` daemon.

Command expose daemon, one-shot execution and client command line interface.

## Command line interface

Command arguments are reported in following sections.

### Help command

You can ask for help al all commands, and the specific execution arguments are explained.

```
go-cron help <command>
```

Available commands are:
* `help`: Help command line arguments
* `explain`: Help explain command line arguments, for accessing to sample input or output formats
* `daemon`: Execute the Scheduler as system daemon, in sync mode
* `once`: Execute the Scheduler as one-shot execution, in sync mode
* `add`: Add a new command in the scheduler and save it to the device
* `remove`: Remove an existing command from the scheduler and save it to the device
* `update`: Update an existing command into the scheduler and save it to the device
* `list`: List command configurations in numerous different output formats
* `active`: List active/running commands in numerous different output formats
* `next`: List next execution of active commands in numerous different output formats


### Explain command

You can ask for input or output samples in numerous different output formats.

```
go-cron explain <command> [-arg0=value0] [-arg1=value1] ...  [-argN=valueN]
```

Available commands are:
* `add`: Add a new command in the scheduler and save it to the device
* `remove`: Remove an existing command from the scheduler and save it to the device
* `update`: Update an existing command into the scheduler and save it to the device
* `list`: List command configurations in numerous output formats
* `active`: List active/running commands in numerous output formats
* `next`: List next execution of active commands in numerous output formats

Optional command arguments [`add`,`remove`,`update`]:
* `in-format` (string) - Encoding input format (text or file) [available: `json`, `xml`, `yaml`]
* `native-in` (bool) - Native GOB input (text or file) encoding format

Optional command arguments [`list`,`active`,`next`]:


### Other Commands

You can execute multiple commands calling appropriate parameters.

```
go-cron <command> [-arg0=value0] [-arg1=value1] ...  [-argN=valueN]
```

Available commands are:
* `daemon`: Execute the Scheduler as system daemon, in sync mode
* `once`: Execute the Scheduler as one-shot execution, in sync mode
* `add`: Add a new command in the scheduler and save it to the device
* `remove`: Remove an existing command from the scheduler and save it to the device
* `update`: Update an existing command into the scheduler and save it to the device
* `list`: List command configurations in numerous different output formats
* `active`: List active/running commands in numerous different output formats
* `next`: List next execution of active commands in numerous different output formats

#### Base command arguments

Base command line arguments are:
* `format` (string) - Encoding file format [available: `json`, `xml`, `yaml`] 
* `path` (string) - Configuration file location 
* `silent` (bool) - Execute less details output for command execution 

#### Daemon command

Executes an asynchronous process in sync mode, accordingly to required base and specific arguments.

Specific command line arguments are:
`no specific argument required` 


#### Once command

Executes a synchronous process, accordingly to required base and specific arguments.

Specific command line arguments are:
`no specific argument required` 


#### Add command

Add a command to the scheduler, accordingly to required base (all mandatory arguments) and specific arguments.

```
go-cron add [-arg0=value0] [-arg1=value1] ...  [-argN=valueN]
```

Specific command line arguments are:
* `in-format` (string) - Encoding input format (text or file) [available: `json`, `xml`, `yaml`]
* `native-in` (bool) - Native GOB input (text or file) encoding format
* `in-file` (string) - Input file absolute path
* `in-text` (sting) - Input text value
* `out-format` (string) - Encoding output format [available: `text`, `json`, `xml`, `yaml`]
* `native-out` (bool) - Native GOB output encoding format


#### Remove command

Remove an existing command in the scheduler at a specific index or in range of indexes, accordingly to required base (all mandatory arguments) and specific arguments. Index is the line number of the `list` command.

```
go-cron remove [-arg0=value0] [-arg1=value1] ...  [-argN=valueN]
```

Specific command line arguments are:
* `index` (int) - Output list raw line number to be deleted (1..n)
* `from` (int) - Output list raw first line number to be deleted (1..n)
* `to` (int) - Output list raw last line number to be deleted (1..n)
* `out-format` (string) - Encoding output format [available: `text`, `json`, `xml`, `yaml`]
* `native-out` (bool) - Native GOB output encoding format


#### Update command

Update an existing command in the scheduler at a specific index, accordingly to required base (all mandatory arguments) and specific arguments. Index is the line number of the `list` command.

```
go-cron update [-arg0=value0] [-arg1=value1] ...  [-argN=valueN]
```

Specific command line arguments are:
* `index` (int) - Output list raw line number to be deleted (1..n)
* `in-format` (string) - Encoding input format (text or file) [available: `json`, `xml`, `yaml`]
* `native-in` (bool) - Native GOB input (text or file) encoding format
* `in-file` (string) - Input file absolute path
* `in-text` (sting) - Input text value
* `out-format` (string) - Encoding output format [available: `text`, `json`, `xml`, `yaml`]
* `native-out` (bool) - Native GOB output encoding format


#### List command

Show the list of saved or cached command configurations with base (all mandatory arguments) and specific arguments, in numerous output encoding format.

```
go-cron list [-arg0=value0] [-arg1=value1] ...  [-argN=valueN]
```

Specific command line arguments are:
* `query` (sting) - Comma separated column_name=value key pairs (*not implemented*)
* `filter` (sting) - Go style template output filter template text (*not implemented*)
* `filter-file` (sting) - Go style template output filter template template (*not implemented*)
* `details` (bool) - Show detailed output format
* `out-format` (string) - Encoding output format [available: `text`, `json`, `xml`, `yaml`]
* `native-out` (bool) - Native GOB output encoding format


#### Active command

Show the list of active/running command execution data with base (all mandatory arguments) and specific arguments, in numerous output encoding format.

```
go-cron active [-arg0=value0] [-arg1=value1] ...  [-argN=valueN]
```

Specific command line arguments are:
* `query` (sting) - Comma separated column_name=value key pairs (*not implemented*)
* `filter` (sting) - Go style template output filter template text (*not implemented*)
* `filter-file` (sting) - Go style template output filter template template (*not implemented*)
* `details` (bool) - Show detailed output format
* `out-format` (string) - Encoding output format [available: `text`, `json`, `xml`, `yaml`]
* `native-out` (bool) - Native GOB output encoding format


#### Next command

Show the list of next running command execution data with base (all mandatory arguments) and specific arguments, in numerous output encoding format.

```
go-cron next [-arg0=value0] [-arg1=value1] ...  [-argN=valueN]
```

Specific command line arguments are:
* `query` (sting) - Comma separated column_name=value key pairs (*not implemented*)
* `filter` (sting) - Go style template output filter template text (*not implemented*)
* `filter-file` (sting) - Go style template output filter template template (*not implemented*)
* `details` (bool) - Show detailed output format
* `out-format` (string) - Encoding output format [available: `text`, `json`, `xml`, `yaml`]
* `native-out` (bool) - Native GOB output encoding format


## DevOps

Installation and build procedures are reported in following sections.



### Build the project

Build command :

```
go build -buildmode=exe github.com/hellgate75/go-cron/
```

### Get the executable

Install locally the command :

```
go get -u github.com/hellgate75/go-cron/
```

### References

Here list of known implementing repositories:

* [Synapses AI](https://github.com/hellgate75/synapses) - Synapse AI open source project

* [Go-Peer-Cluster](https://github.com/hellgate75/go-peer-nodes) - Peer-to-Peer auto-gossip and auto-grouping Cluster-ware library

Enjoy the experience.


## License

The library is licensed with [LGPL v. 3.0](/LICENSE) clauses, with prior authorization of author before any production or commercial use. Use of this library or any extension is prohibited due to high risk of damages due to improper use. No warranty is provided for improper or unauthorized use of this library or any implementation.

Any request can be prompted to the author [Fabrizio Torelli](https://www.linkedin.com/in/fabriziotorelli) at the following email address:

[hellgate75@gmail.com](mailto:hellgate75@gmail.com)
 
