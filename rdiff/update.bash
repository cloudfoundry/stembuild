# This is only here as an example and has not been tested!
#
# This *must* be run from the rdiff directory where this
# file is located!

# Use this fork of librsync: https://github.com/charlievieth/librsync/tree/mingw64-fseeko64-v2.0.0
#
# Root directory of librsync
: "${LIBRSYNC_DIR:?Need to set LIBRSYNC_DIR non-empty}"

if [[ ! -d "$LIBRSYNC_DIR" ]]; then
	echo "Invalid LIBRSYNC_DIR: $LIBRSYNC_DIR"
fi

# Directory where generated Makefile and generate source code will
# be saved.
LIBRSYNC_BUILD_DIR="$LIBRSYNC_DIR/build"

mkdir "$LIBRSYNC_BUILD_DIR"

if [[ ! -d "$LIBRSYNC_BUILD_DIR" ]]; then
	echo "Invalid LIBRSYNC_BUILD_DIR: $LIBRSYNC_BUILD_DIR"
fi

# # This is where the library would be saved if 'make' is ran on the
# # resulting build files.  This is not actually used, but is good to
# # set (in case someone decides to run 'make install') as it prevents
# # overwriting the librsync library in /usr/local.
LIBRSYNC_INSTALL_DIR="$LIBRSYNC_DIR/install"

mkdir "$LIBRSYNC_INSTALL_DIR"

# Run cmake
pushd "$LIBRSYNC_BUILD_DIR"
	cmake \
	  -DCMAKE_INSTALL_PREFIX:PATH="$LIBRSYNC_INSTALL_DIR" \
	  -DCMAKE_BUILD_TYPE:STRING="Release" \
	  "$LIBRSYNC_DIR"
popd

# Remove any C source code artifacts from rdiff directory.
rm -vf *.[ch]

# We need glob expansion - so $LIBRSYNC_DIR better not contain
# any spaces!!
cp $LIBRSYNC_BUILD_DIR/src/*.[ch] .
cp $LIBRSYNC_DIR/src/*.[ch] .

rm -v rdiff.c
