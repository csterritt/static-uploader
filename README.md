### static-uploader

This program takes a directory and a bucket (in an S3-compatible storage facility)
and uploads the files below that directory to the bucket. Any files that are found
in the bucket whose name does not match a file under the directory is removed. Any
files that are the same name from the last upload are replaced with the (possibly
new) file.

The directory is descended recursively, so all files are found, and are uploaded with
the proper subdirectory prefix.

This is useful for, e.g., `dist` directories for web applications, where the bucket
holds the files for a website.
