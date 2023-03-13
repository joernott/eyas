EYaS - EYaml Server - Documentation
-----------------------------------

[Purpose](#docu_purpose)
[Usage](#docu_usage)
[Frontend](#docu_frontend)
[Installation](#docu_installation)
[Configuration](#docu_configuration)

### Purpose

EYAaS is a simple web frontend to [eyaml](https://github.com/voxpupuli/hiera-eyaml) encryption. Its intended audience are users who don't have eyaml installed or access to the public keys as well as lazy system administrators (like me) who need to generate/encrypt a lot of passwords at once.

Written in go and using the original eyaml cli as the main workhorse, EYas provides a simple interface to encrypt a single password, a whole yaml- or csv file for one or more puppet servers (using different pkcs7 keys).

It aims for easy setup and configuration as well as convenience.

### Usage

#### systemd service

When using the provided RPM, a [systemd](https://systemd.io/) [service](https://www.freedesktop.org/software/systemd/man/systemd.service.html#) will be provided. Use systemctl start eyas.service to start it and systemctl enable eyas.service to start it automatically when the server boots.

#### docker / podman container

When using the provided Dockerfile to build a container, this container can be started with podman run --rm --expose=8443 --volume=./keys:/keys:ro --volume=./ssl:/ssl:ro eyas:latest.

#### Start manually

When starting EYaS without parameters, it will start a web server using its default settings.

The web server uses https by default and listens on port 8443. Specify --ssl=false to use http instead. This is insecure, as passwords can be transmitted as clear text and should only be used when you can secure communications between the web browser and the server by other means (reverse proxy).

EYaS assumes to find its ssl certificate and key in the "ssl" directory in its working directory named server.crt and server.key. You can use [Let's Encrypt](https://letsencrypt.org/de/) to create and manage these certificates or create your own. Currently, golang does not support encrypted keys, so the key needs to be unencrypted. Refer to [the setup documentation](#setup) for more information.

It will assume, that "eyaml" is somewhere in its path and the directory "keys" contains folders for every puppetserver which in turn contain a "public\_key.pkcs7.pem" with the servers public key. You can use the --keydir parameter to choose a diffent directory.

The application logs tom stdout by default. This can be changed by providing the name of a log file with the --logfile parameter. The log level can be set with the --loglevel parameter. Valid log levels are panic, fatal, error, warn, info, debug and trace. By default, the log level is set to "info".

### Web Frontend

EYaS is quite simple to use. The menu on the left side of the screen has three buttons "Single key", "YAML" and "CSV" where you can choose whether you want to encrypt a single password or a whole list in either YAML or CSV format. The button "Output" shows the last output and when you read this, chances are quite high, that you already figured out, what the button "Documentation" is for.

#### Encrypting a single key

To encrypt a single sensitive key, you select the page by clicking on the "Single Key" entry. The name of your key goes into the "Key" field. If you want to generate an entry for the parameter "mypassword" of the puppet module "mymodule", you write "mymodule::mypassword" here.

The fields "Password" and "Repeat password" are used to set the password. They must match.

Use "Output type" to specify whether you want the encrypted result as a multiline block or a single very long line of text.

The list in "Use PKCS7 key(s)"" shows all available public keys which can be used for encryption. If you need to set up the same entry across multiple puppet servers which use different key pairs, you can select multiple keys at once to generate them together.

Press "Encrypt" to encrypt the key. The content of your form will be processed by the server and upon success, you will be redirected to the "Output" page

#### Encrypt a yaml file

Use the "YAML" button in the menu to encrypt a hiera yaml file in one go. Paste the content into the "YAML" textarea, choose the output type and PKCS7 keys as described above and click on "Encrypt". Depending on the number of entries, this may take a while. The output page will contain the output for every selected key.

#### Encrypt passwords from a CSV list

Several password managers can export their data as comma separated text (sometimes with other separators). Use the "CSV" button to get to the form. Here, you paste the contents of your file into the CSV text area. Then you specify, in which column the key and in which the password can be found, by typing their number into the "Key column" and "Password column" fields. Column numbering starts at 0. The field "Separator" specifies, which separator is used between the fields (defaults to a comma).

"Output type" and "Use PKCS7 key(s)" work the same way as described above for the single key. Klick on "Encrypt" to find the results on the "Output" page. Depending on the number of entries, the encryption may take a while.

### Installation

#### RPM

On a RHEL compatible system, use the RPM provided on the [release page](https://github.com/joernott/eyas/releases) to install EYaS into /usr/lib/eyas. The RPM will also provide a systemd service to start the server.

Download the rpm from the [release page](https://github.com/joernott/eyas/releases) and save it somewhere (e.g. tmp). Then use your system installer to install it. Use the [instructions on the hiera-eyaml Site](https://github.com/voxpupuli/hiera-eyaml#setup) to install eyaml.

On Redhat 7 / CentOS 7:

`sudo yum install /path/to/eyas/rpm sudo gem install hiera-eyaml`

On Redhat 8 / Fedora etc

`sudo dnf install /path/to/eyas/rpm sudo gem install hiera-eyaml`

#### Other Linux distros

Use the tar.gz archive provided on the [release page](https://github.com/joernott/eyas/releases) to install eyas. Afterwards copy the systemd service file to the right place and install eyaml.

`sudo mkdir /usr/lib/eyas sudo tar -C /usr/lib/eyas -xzf /path/to/eyas/archive sudo mv /usr/lib/eyas/eyas.service /usr/lib/systemd/system/ sudo systemctl service-reload sudo gem install hiera-eyaml`

#### Windows

Figure out, how to install ruby and hiera eyaml on Windows. Then create a directory "eyas" in "C:\\Program Files" and copy the eyas.exe from the [release page](https://github.com/joernott/eyas/releases) into that folder.

### Configuration

#### Commandline

EYAS uses cobra/viper for parsing the commandline flags and the config file "eyas.yaml". The following parameters are supported:

\--config string

Config file (default "eyas.yaml")

\-h, --help

help for eyas

\-k, --keydir string

Folder containing directories with PKCS7 public keys (default "keys")

\-l, --logfile string

Where to log, defaults to stdout

\-L, --loglevel string

Loglevel can be one of panic, fatal, error, warn, info, debug, trace (default "info")

\-p, --port int

Port to run the server on (default 8443)

\-s, --ssl

Whether to use SSL or not (default true)

#### Keys

For every PKCS7 key you want to encrypt for, the keydir must contain a folder which in turn contains a file named "public\_key.pkcs7.pem" containing the PUBLIC key to encrypt for. The name of the folders are used in the UI to select the keys.