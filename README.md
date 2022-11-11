# Weak Table

In an effort to bring the concept of weak tables to Go this project was created to make a way for transient data to be stored against a weak reference key.  The main value for this is to add metadata to an already existing `struct` pointer you have no control over.  This way you can extend it to store extra data it didn't previously support.

This already has uses in cases where you want to inject extra metadata on an object such as for dependency injection control or in libraries that so not support Go fully or provide only `struct` definitions when an interface would give you the flexibility you wanted.