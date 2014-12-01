Hydra
=====

Framework to create distributed calculations.

This framework is in the early stages and currently does not yet support distribution.  The framework links several functions together and can create several instances of each function.

How It Works
============

Functions decide which of the three categories they fall under:

1. Producer
2. Filter
3. Consumer

Producer
========

A ```Producer``` does not have an input but has an output.  It is intended to begin the pipeline of data.

Filter
======

A ```Filter``` has an input and output.  It is intended to manipulate the data and pick and choose what data is passed to the next function.

Consumer
========

A ```Consumer``` has an input but no output.  This is a final stage in the pipeline.

Mechanics
=========

Functions communicate through channels.
