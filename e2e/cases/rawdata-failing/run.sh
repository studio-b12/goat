goat body_fail.goat && {
    echo "should have failed but didn't"
    exit 1
}

goat formdata_fail.goat && {
    echo "should have failed but didn't"
    exit 1
}

exit 0