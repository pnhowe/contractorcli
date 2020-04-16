contractorcli
=============

contractorcli is a commandline utility (simmaler to docker/kubectl/esxcli)
that you run on your local machine to talk to an instance of contractor.

See https://t3kton.github.io/contractor/contractorcli.html for more information

Pre-Build Binaries
------------------

You can get a pre-built binary from https://github.com/T3kton/contractorcli/releases

Building from Source
--------------------

You will need make and go of at least version 1.12 (building uses the new module system)

clone the source::

  git clone https://github.com/T3kton/contractorcli.git

cd into the source::

  cd contractorcli

and compile::

  make

go will download the required dependencies and build contractorcli.

Enjoy!


